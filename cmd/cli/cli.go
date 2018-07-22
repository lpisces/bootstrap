package cli

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/labstack/gommon/log"
	ccli "gopkg.in/urfave/cli.v1"
	"os"
	"sync"
)

const (
	MaxClient = 100
)

func Run(c *ccli.Context) (err error) {

	// config file
	config := c.String("config")
	if config != "" {
		if err = Conf.LoadFromIni(config); err != nil {
			return
		}
	}

	// env
	env := c.String("env")
	if env != "" {
		Conf.Env = env
	}

	// debug
	if c.Bool("debug") {
		Conf.Debug = true
	}

	// source
	if c.String("source") != "" {
		Conf.Source = c.String("source")
	}

	// output
	if c.String("output") != "" {
		Conf.Output = c.String("output")
	}

	// cache
	if c.String("cache") != "" {
		Conf.CacheDir = c.String("cache")
	}
	if _, err := os.Stat(Conf.CacheDir); os.IsNotExist(err) {
		if err = os.Mkdir(Conf.CacheDir, os.ModeDir); err != nil {
			return err
		}
	}

	// kuaidi100 customer id
	if c.String("customer") != "" {
		Conf.Kuaidi100Config.Customer = c.String("customer")
	}

	// kuaidi100 key
	if c.String("key") != "" {
		Conf.Kuaidi100Config.Key = c.String("key")
	}

	f, err := os.Open(Conf.Source)
	if err != nil {
		return
	}

	r := csv.NewReader(bufio.NewReader(f))
	source, _ := r.ReadAll()

	packages := make(chan *Package, len(source))
	var wg sync.WaitGroup
	pChan := make(chan int, MaxClient)
	for i := range source {
		if len(source[i]) < 3 {
			continue
		}
		pChan <- i
		wg.Add(1)
		go func() {
			p := &Package{
				source[i][0],
				source[i][1],
				source[i][2],
				false,
				false,
				"postnl",
				0,
			}
			if p.GlobalExpressCode != "" {
				p.TraceCN()
				p.TraceNL()
			}
			packages <- p
			defer func() {
				wg.Done()
				<-pChan
			}()
		}()
	}

	wg.Wait()
	close(packages)

	fd, err := os.Create(Conf.Output)
	if err != nil {
		return err
	}
	defer fd.Close()

	fd.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(fd)

	for p := range packages {
		if Conf.Debug {
			log.Info(p)
		}
		w.Write([]string{p.GlobalExpressCode, p.CNExpressCode, p.CNPostCode, fmt.Sprintf("%v", p.CNStatus), fmt.Sprintf("%v", p.NLStatus), fmt.Sprintf("%v", p.Weight)})
	}
	w.Flush()

	return
}
