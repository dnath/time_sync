Time Synchronization : Marzullo's Algorithm
===========================================


A Go client for syncing local time with multiple time servers according to Marzullo's algorithm.

For running it, you have to have [golang](http://golang.org/) installed on your machine. Time server ip's have been pre-configured in the client.

Running it is very simple, just type the following into the terminal and viola...

	$ go run client.go
	
Output will of the form --
	
	Trial #0:
	current time = 2014-10-24 17:26:21.926508032 -0700 PDT
	[54.172.168.244:5000] server time = 2014-10-24 17:26:21.93695744 -0700 PDT, rtt = 0.077063s
	[54.169.67.45:5000] server time = 2014-10-24 17:26:21.988278272 -0700 PDT, rtt = 0.181922s
	[54.207.15.207:5000] server time = 2014-10-24 17:26:22.005858304 -0700 PDT, rtt = 0.193944s
	[54.191.73.92:5000] server time = 2014-10-24 17:26:21.91238272 -0700 PDT, rtt = 0.037410s
	[128.111.44.106:12291] server time = 2014-10-24 17:26:23.33696256 -0700 PDT, rtt = 0.002098s
	updated time = 2014-10-24 17:26:21.91238272 -0700 PDT, error = 0.037410s

	Trial #1:
	current time = 2014-10-24 17:26:22.419638016 -0700 PDT
	[54.172.168.244:5000] server time = 2014-10-24 17:26:22.42916352 -0700 PDT, rtt = 0.071990s
	[54.169.67.45:5000] server time = 2014-10-24 17:26:22.492625152 -0700 PDT, rtt = 0.193896s
	[54.207.15.207:5000] server time = 2014-10-24 17:26:22.506194176 -0700 PDT, rtt = 0.201283s
	[54.191.73.92:5000] server time = 2014-10-24 17:26:22.405045248 -0700 PDT, rtt = 0.037312s
	[128.111.44.106:12291] server time = 2014-10-24 17:26:23.829936128 -0700 PDT, rtt = 0.001856s
	updated time = 2014-10-24 17:26:22.405045248 -0700 PDT, error = 0.037312s

	Trial #2:
	current time = 2014-10-24 17:26:22.926302976 -0700 PDT
	[54.172.168.244:5000] server time = 2014-10-24 17:26:22.935799808 -0700 PDT, rtt = 0.072087s
	[54.169.67.45:5000] server time = 2014-10-24 17:26:22.998688512 -0700 PDT, rtt = 0.201000s
	[54.207.15.207:5000] server time = 2014-10-24 17:26:23.029405952 -0700 PDT, rtt = 0.253518s
	[54.191.73.92:5000] server time = 2014-10-24 17:26:22.912174336 -0700 PDT, rtt = 0.036454s
	[128.111.44.106:12291] server time = 2014-10-24 17:26:24.336781568 -0700 PDT, rtt = 0.002276s
	updated time = 2014-10-24 17:26:22.912174336 -0700 PDT, error = 0.036454s
	


