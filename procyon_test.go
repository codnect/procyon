package procyon

import (
	"testing"
)

func init() {
	/*procyonAppFlagSet := flag.NewFlagSet("procyon", flag.ContinueOnError)
	procyonAppFlagSet.Parse(os.Args)

	flag.Visit(func(f *flag.Flag) {
		log.Printf("")
	})
	flag.Bool("fork", false, "hey")
	flag.NewFlagSet("", flag.ContinueOnError)
	flag.Parse()
	log.Print()*/
}

func TestProcyonApplication(t *testing.T) {
	app := NewProcyonApplication()
	app.SetApplicationRunListeners()
	app.Run()
	//assert.Equal(t, true, true)
}
