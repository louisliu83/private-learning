package config

const (
	defaultControllerListen  string = "0.0.0.0:8080"
	defaultStreamListen      string = "0.0.0.0:7766"
	defaultDatasetValidDays  int32  = 30
	defaultDatasetMaxLine    int64  = 1000 * 10000
	defaultServerConcurrency int    = 1
	defaultClientConcurrency int    = 1
)

func (c *Config) setDefaultConfigIfNotSet() {
	if len(c.Listener.Address) == 0 {
		c.Listener.Address = defaultControllerListen
	}
	if len(c.StreamListener.Address) == 0 {
		c.StreamListener.Address = defaultStreamListen
	}
	if c.DataSet.ValidDays == 0 {
		c.DataSet.ValidDays = defaultDatasetValidDays
	}
	if c.DataSet.MaxLines == 0 {
		c.DataSet.MaxLines = defaultDatasetMaxLine
	}
	if c.DataSet.Sharders == 0 {
		c.DataSet.Sharders = 1
	}
	if c.DataSet.Downloaders == 0 {
		c.DataSet.Downloaders = 1
	}
	if c.PsiExecutor.ServerConcurrency == 0 {
		c.PsiExecutor.ServerConcurrency = defaultServerConcurrency
	}
	if c.PsiExecutor.ClientConcurrency == 0 {
		c.PsiExecutor.ClientConcurrency = defaultClientConcurrency
	}
}
