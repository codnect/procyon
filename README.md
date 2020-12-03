
<img src="https://procyon-projects.github.io/img/logo.png" width="128">

# Procyon Framework
![alt text](https://goreportcard.com/badge/github.com/procyon-projects/procyon)


# What is Procyon? 

Procyon is a HTTP web framework written in Go, powered by [fasthttp](https://github.com/valyala/fasthttp) 
and third-part libraries. If you need a comprehensive web framework, then use Procyon.
Because it provides a lot of modules which include several features.

* It makes it easy to create production-grade applications. 
* It aims to ease to build, develop and deploy your web applications quickly in Go.

## Third-Party Libraries

We use some third-party libraries while developing Procyon. Here are the list we've used :

* [fasthttp](https://github.com/valyala/fasthttp), It is 10 times faster than standard http library. 
* [jsoniter](https://github.com/json-iterator/go), It is used to encode and decode. It's faster than standard library.

## Modules
There are a number of modules in Procyon Framework :

* **procyon :**  It provides all features of Procyon.

* **procyon-core :** This provides core features and utilities for all modules.

* **procyon-configure :** This includes and provides configurations that configure application automatically.

* **procyon-context :** This provides context for Procyon applications.

* **procyon-web :** It provides web support for developing web application.

* **procyon-peas :** This allow us to manage our instances created in Procyon application. Peas are very similar to Java Beans.
They might be called as Go Beans :)


## How to use Procyon?

It is so easy to use Procyon Framework. The only thing you have to do is to add the **Procyon** module into your **go.mod** and 
import it into your code file.
```go
import (
	"github.com/procyon-projects/procyon"
)
```
Next, You need to invoke the method **procyon.NewProcyonApplication** to create a procyon application in main function.
```go
myApp := procyon.NewProcyonApplication()
```
After that, invoke the method **Run** to run the application. It is easy that much to have a simple Procyon application.
```go
myApp.Run()
``` 
Eventually, your code snippet will look like the following.  It is easy that much to have a simple Procyon application.
```go
import (
	"github.com/procyon-projects/procyon"
)

func main() {
	myApp := procyon.NewProcyonApplication()
	myApp.Run()
}
```

After running, you will see the following console.

![Image of Yaktocat](https://procyon-projects.github.io/img/run-console.png)

# QuickStart
This quickstart gives you a basic understanding of creating a simple endpoint and how to do
by using Procyon Framework.

## Components
Components are the part of Procyon Framework. Controller, Service, Repository, Initializers, etc.
are considered as component.  

## Register Components
If you have a component like a Controller, you need to notify it to Procyon as Go language
doesn't have annotation and similar reflection library like Java. You don't need to be worried about it.
The only thing you have to do is to use method **core.Register** in **init**, as the following below.
It is placed in module **procyon-core**.

**Note that you have to give a function with only one return parameter, which will create a instance of the controller.**

**config.go**
```go
import (
	core "github.com/procyon-projects/procyon-core"
	"github.com/procyon-projects/procyon-test-app/controller"
)

func init() {
	/* Controllers */
	core.Register(controller.NewHelloWorldController)
}

```
## First Controller
The first thing you have to do to have a controller component in Procyon is to implement
interface **Controller**.

The interface **Controller** looks like the following below. It is placed in module **procyon-web**.

```go
type Controller interface {
	RegisterHandlers(registry HandlerRegistry)
}
```

A simple controller will looks like the following below. It might change based on your needs.

```go
type HelloWorldController struct {
}

func NewHelloWorldController() HelloWorldController {
	return HelloWorldController{}
}

func (controller HelloWorldController) RegisterHandlers(registry web.HandlerRegistry) {
	...
}

func (controller HelloWorldController) HelloWorld(context *web.WebRequestContext) {
    context.Ok().SetBody("Hello World")
}
```

## Registry Handlers
Your handler registrations should be done by using **registry** which will 
be passed into the method **RegisterHandlers**, as you can see the following.

```go
func (controller HelloWorldController) RegisterHandlers(registry web.HandlerRegistry) {
	registry.Register(
		web.Get(controller.HelloWorld, web.Path("/api/helloworld"))
	)
}
```

You can see the complete code below. It is easy that much to create a controller and register
your handlers.

**controller.go**
```go
type HelloWorldController struct {
}

func NewHelloWorldController() NewHelloWorldController {
	return HelloWorldController{}
}

func (controller HelloWorldController) RegisterHandlers(registry web.HandlerRegistry) {
	registry.Register(
		web.Get(controller.HelloWorld, web.Path("/api/helloworld"))
	)
}

func (controller HelloWorldController) HelloWorld(context *web.WebRequestContext) {
    context.Ok().SetBody("Hello World")
}
```

## Run Application
If you run your application without giving any parameter port, it will start on port 8080
as you can see following.

![Image of Console Running](https://procyon-projects.github.io/img/run-console.png)

If you want to change the port on which the application start, you can specify parameter **--server.port**.
When you specify the parameter port as 3030 (**--server.port=3030**), your application will start on **port 3030**. 


## Request the endpoint
We assume that you do all steps which needs to be done and have a running application, then it's time to request
the endpoint **/api/helloworld** which we create.

![Image of Console Running](https://procyon-projects.github.io/img/api-helloworld.png)

# Benchmarks
You can find the benchmark result including memory consumption below and compare Procyon with other router and frameworks.
We got the benchmark results by using [go-http-routing-benchmark](https://github.com/procyon-projects/go-http-routing-benchmark).

## Benchmark System

* Intel Core i7-4700MQ 2.40GHz
* 8 GiB RAM
* go version 1.13.4 windows/amd64

## Memory Consumption

##### GITHUB API ROUTES : 203
**Procyon:** 108424 Bytes
```
#GithubAPI Routes: 203
   Ace: 48640 Bytes
   Aero: 235280 Bytes
   Bear: 82328 Bytes
   Beego: 150936 Bytes
   Bone: 100976 Bytes
   Chi: 95112 Bytes
   CloudyKitRouter: 93704 Bytes
   Denco: 36448 Bytes
   Echo: 100040 Bytes
   Gin: 58512 Bytes
   GocraftWeb: 95640 Bytes
   Goji: 49680 Bytes
   Gojiv2: 104704 Bytes
   GoJsonRest: 142104 Bytes
   GoRestful: 1241656 Bytes
   GorillaMux: 1322784 Bytes
   GowwwRouter: 80008 Bytes
   HttpRouter: 37096 Bytes
   HttpTreeMux: 78800 Bytes
   Kocha: 785552 Bytes
   LARS: 48600 Bytes
   Macaron: 92784 Bytes
   Martini: 485264 Bytes
   Pat: 21200 Bytes
   Possum: 85600 Bytes
   Procyon: 108424 Bytes
   R2router: 47104 Bytes
   Rivet: 42840 Bytes
   Tango: 54840 Bytes
   TigerTonic: 95248 Bytes
   Traffic: 921760 Bytes
   Vulcan: 425368 Bytes
```

##### GPLUS API ROUTES : 203
**Procyon:** 9752 Bytes
```
#GPlusAPI Routes: 13
   Ace: 3664 Bytes
   Aero: 26552 Bytes
   Bear: 7112 Bytes
   Beego: 10272 Bytes
   Bone: 6688 Bytes
   Chi: 8024 Bytes
   CloudyKitRouter: 6728 Bytes
   Denco: 3264 Bytes
   Echo: 9640 Bytes
   Gin: 4384 Bytes
   GocraftWeb: 7496 Bytes
   Goji: 3152 Bytes
   Gojiv2: 7376 Bytes
   GoJsonRest: 11416 Bytes
   GoRestful: 74328 Bytes
   GorillaMux: 66208 Bytes
   GowwwRouter: 5744 Bytes
   HttpRouter: 2760 Bytes
   HttpTreeMux: 7440 Bytes
   Kocha: 128880 Bytes
   LARS: 3656 Bytes
   Macaron: 8656 Bytes
   Martini: 23920 Bytes
   Pat: 1856 Bytes
   Possum: 7248 Bytes
   Procyon: 9752 Bytes
   R2router: 3928 Bytes
   Rivet: 3064 Bytes
   Tango: 5168 Bytes
   TigerTonic: 9408 Bytes
   Traffic: 46400 Bytes
   Vulcan: 25752 Bytes
```


##### PARSE API ROUTES : 203
**Procyon:** 15024 Bytes
```
#ParseAPI Routes: 26
   Ace: 6656 Bytes
   Aero: 29304 Bytes
   Bear: 12320 Bytes
   Beego: 19280 Bytes
   Bone: 11440 Bytes
   Chi: 9744 Bytes
   CloudyKitRouter: 11208 Bytes
   Denco: 4192 Bytes
   Echo: 11824 Bytes
   Gin: 7776 Bytes
   GocraftWeb: 12800 Bytes
   Goji: 5680 Bytes
   Gojiv2: 14464 Bytes
   GoJsonRest: 14072 Bytes
   GoRestful: 116264 Bytes
   GorillaMux: 105880 Bytes
   GowwwRouter: 9344 Bytes
   HttpRouter: 5024 Bytes
   HttpTreeMux: 7848 Bytes
   Kocha: 181712 Bytes
   LARS: 6632 Bytes
   Macaron: 13648 Bytes
   Martini: 45888 Bytes
   Pat: 2560 Bytes
   Possum: 9200 Bytes
   Procyon: 15024 Bytes
   R2router: 7056 Bytes
   Rivet: 5680 Bytes
   Tango: 8920 Bytes
   TigerTonic: 9840 Bytes
   Traffic: 79096 Bytes
   Vulcan: 44504 Bytes
```

##### STATIC ROUTES : 203
**Procyon:** 37976 Bytes
```
#Static Routes: 157
   HttpServeMux: 14512 Bytes
   Ace: 30648 Bytes
   Aero: 34536 Bytes
   Bear: 31080 Bytes
   Beego: 98456 Bytes
   Bone: 40224 Bytes
   Chi: 83608 Bytes
   CloudyKitRouter: 30448 Bytes
   Denco: 9928 Bytes
   Echo: 80280 Bytes
   Gin: 34936 Bytes
   GocraftWeb: 55496 Bytes
   Goji: 29744 Bytes
   Gojiv2: 105840 Bytes
   GoJsonRest: 138872 Bytes
   GoRestful: 816936 Bytes
   GorillaMux: 585632 Bytes
   GowwwRouter: 24968 Bytes
   HttpRouter: 21680 Bytes
   HttpTreeMux: 73448 Bytes
   Kocha: 115472 Bytes
   LARS: 30640 Bytes
   Macaron: 38592 Bytes
   Martini: 310864 Bytes
   Pat: 19696 Bytes
   Possum: 89920 Bytes
   Procyon: 37976 Bytes
   R2router: 23712 Bytes
   Rivet: 24608 Bytes
   Tango: 28264 Bytes
   TigerTonic: 78560 Bytes
   Traffic: 538976 Bytes
   Vulcan: 369960 Bytes
```

## Micro Benchmarks

```
BenchmarkProcyon_Param                	 3831451	       313 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_Param5               	 2949872	       379 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_Param20              	 4731327	       243 ns/op	       0 B/op	       0 allocs/op
```
```
BenchmarkAce_Param                    	 1872678	       721 ns/op	      32 B/op	       1 allocs/op
BenchmarkAero_Param                   	 5191189	       232 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param                   	  600434	      2923 ns/op	     456 B/op	       5 allocs/op
BenchmarkBeego_Param                  	  299784	      5864 ns/op	     416 B/op	       7 allocs/op
BenchmarkBone_Param                   	  444094	      5528 ns/op	     816 B/op	       6 allocs/op
BenchmarkChi_Param                    	  706322	      3401 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_Param        	 9829281	       145 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_Param                  	 2085495	       513 ns/op	      32 B/op	       1 allocs/op
BenchmarkEcho_Param                   	 3667047	       343 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param                    	 3190291	       395 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param             	  444615	      4236 ns/op	     648 B/op	       8 allocs/op
BenchmarkGoji_Param                   	  775995	      1963 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Param                 	  211291	      8051 ns/op	    1328 B/op	      11 allocs/op
BenchmarkGoJsonRest_Param             	  363358	      5418 ns/op	     649 B/op	      13 allocs/op
BenchmarkGoRestful_Param              	   57375	     20645 ns/op	    4192 B/op	      14 allocs/op
BenchmarkGorillaMux_Param             	  119913	     10017 ns/op	    1280 B/op	      10 allocs/op
BenchmarkGowwwRouter_Param            	 1000000	      2783 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_Param             	 2519222	       425 ns/op	      32 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Param            	  997340	      2110 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_Param                  	 1000000	      1093 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_Param                   	 4178226	       288 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param                	  196654	      8871 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_Param                	   63776	     17606 ns/op	    1072 B/op	      10 allocs/op
BenchmarkPat_Param                    	  333085	      4617 ns/op	     536 B/op	      11 allocs/op
BenchmarkPossum_Param                 	  749427	      4251 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_Param                	 3831451	       313 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_Param               	  799248	      2508 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Param                  	 1525638	       747 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_Param                  	  413841	      4217 ns/op	     248 B/op	       8 allocs/op
BenchmarkTigerTonic_Param             	  210378	      7398 ns/op	     776 B/op	      16 allocs/op
BenchmarkTraffic_Param                	   80450	     12800 ns/op	    1856 B/op	      21 allocs/op
BenchmarkVulcan_Param                 	  999291	      2123 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_Param5                   	 1000000	      1673 ns/op	     160 B/op	       1 allocs/op
BenchmarkAero_Param5                  	 4466318	       323 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param5                  	  413829	      3773 ns/op	     501 B/op	       5 allocs/op
BenchmarkBeego_Param5                 	  239832	      7152 ns/op	     480 B/op	       7 allocs/op
BenchmarkBone_Param5                  	  299614	      5905 ns/op	     864 B/op	       6 allocs/op
BenchmarkChi_Param5                   	  461614	      4607 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_Param5       	 2118181	       593 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_Param5                 	 1000000	      1382 ns/op	     160 B/op	       1 allocs/op
BenchmarkEcho_Param5                  	 1564947	       763 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param5                   	 1710350	       692 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param5            	  260689	      7709 ns/op	     920 B/op	      11 allocs/op
BenchmarkGoji_Param5                  	  704568	      3151 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Param5                	  118723	      9216 ns/op	    1392 B/op	      11 allocs/op
BenchmarkGoJsonRest_Param5            	  130878	     10705 ns/op	    1097 B/op	      16 allocs/op
BenchmarkGoRestful_Param5             	   52820	     22453 ns/op	    4288 B/op	      14 allocs/op
BenchmarkGorillaMux_Param5            	   95931	     14608 ns/op	    1344 B/op	      10 allocs/op
BenchmarkGowwwRouter_Param5           	  922459	      2864 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_Param5            	 1000000	      1256 ns/op	     160 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Param5           	  352465	      5268 ns/op	     576 B/op	       6 allocs/op
BenchmarkKocha_Param5                 	  386816	      3933 ns/op	     440 B/op	      10 allocs/op
BenchmarkLARS_Param5                  	 2445212	       484 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param5               	  128941	      9963 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_Param5               	   83846	     20491 ns/op	    1232 B/op	      11 allocs/op
BenchmarkPat_Param5                   	  149848	     11648 ns/op	     888 B/op	      29 allocs/op
BenchmarkPossum_Param5                	  461324	      4174 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_Param5               	 2949872	       379 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_Param5              	  599678	      3125 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Param5                 	 1000000	      1893 ns/op	     240 B/op	       1 allocs/op
BenchmarkTango_Param5                 	  352447	      5223 ns/op	     360 B/op	       8 allocs/op
BenchmarkTigerTonic_Param5            	   42673	     29858 ns/op	    2279 B/op	      39 allocs/op
BenchmarkTraffic_Param5               	   55767	     21490 ns/op	    2208 B/op	      27 allocs/op
BenchmarkVulcan_Param5                	  444529	      3066 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_Param20                  	  706201	      4108 ns/op	     640 B/op	       1 allocs/op
BenchmarkAero_Param20                 	 1309110	       890 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param20                 	   74020	     15101 ns/op	    1664 B/op	       5 allocs/op
BenchmarkBeego_Param20                	   92810	     15628 ns/op	     544 B/op	       7 allocs/op
BenchmarkBone_Param20                 	   66982	     17958 ns/op	    2030 B/op	       6 allocs/op
BenchmarkChi_Param20                  	  214224	      8188 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_Param20      	  499647	      2342 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_Param20                	  460825	      4528 ns/op	     640 B/op	       1 allocs/op
BenchmarkEcho_Param20                 	  999291	      1783 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param20                  	  705255	      1801 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param20           	   34161	     30874 ns/op	    3795 B/op	      15 allocs/op
BenchmarkGoji_Param20                 	  148042	     10547 ns/op	    1247 B/op	       2 allocs/op
BenchmarkGojiv2_Param20               	   85652	     13179 ns/op	    1632 B/op	      11 allocs/op
BenchmarkGoJsonRest_Param20           	   32763	     40715 ns/op	    4486 B/op	      20 allocs/op
BenchmarkGoRestful_Param20            	   23722	     45437 ns/op	    6716 B/op	      18 allocs/op
BenchmarkGorillaMux_Param20           	   33889	     33634 ns/op	    3451 B/op	      12 allocs/op
BenchmarkGowwwRouter_Param20          	  427831	      3928 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_Param20           	  666177	      3970 ns/op	     640 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Param20          	   40104	     28046 ns/op	    3195 B/op	      10 allocs/op
BenchmarkKocha_Param20                	   69314	     15448 ns/op	    1808 B/op	      27 allocs/op
BenchmarkLARS_Param20                 	 1577797	       749 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param20              	   61416	     17013 ns/op	    2923 B/op	      12 allocs/op
BenchmarkMartini_Param20              	   50074	     22807 ns/op	    3595 B/op	      13 allocs/op
BenchmarkPat_Param20                  	   37208	     30413 ns/op	    4423 B/op	      93 allocs/op
BenchmarkPossum_Param20               	  999433	      2157 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_Param20              	 4731327	       243 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_Param20             	  164264	      8246 ns/op	    2282 B/op	       7 allocs/op
BenchmarkRivet_Param20                	  413500	      3155 ns/op	    1024 B/op	       1 allocs/op
BenchmarkTango_Param20                	  352708	      7856 ns/op	     856 B/op	       8 allocs/op
BenchmarkTigerTonic_Param20           	   10000	    131285 ns/op	    9870 B/op	     119 allocs/op
BenchmarkTraffic_Param20              	   16506	     61298 ns/op	    7852 B/op	      47 allocs/op
BenchmarkVulcan_Param20               	  307441	      3576 ns/op	      98 B/op	       3 allocs/op
```

## Github API Benchmarks
```
BenchmarkProcyon_GithubStatic         	 5793034	       336 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GithubParam          	 2282331	       628 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GithubAll            	    9224	    156204 ns/op	       0 B/op	       0 allocs/op
```
```
BenchmarkAce_GithubStatic             	 2241300	       535 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubStatic            	 4914578	       244 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubStatic            	  799477	      1846 ns/op	     120 B/op	       3 allocs/op
BenchmarkBeego_GithubStatic           	  285523	      4635 ns/op	     416 B/op	       7 allocs/op
BenchmarkBone_GithubStatic            	   36447	     47048 ns/op	    2880 B/op	      60 allocs/op
BenchmarkCloudyKitRouter_GithubStatic 	 4355858	       277 ns/op	       0 B/op	       0 allocs/op
BenchmarkChi_GithubStatic             	  667240	      3417 ns/op	     432 B/op	       3 allocs/op
BenchmarkDenco_GithubStatic           	 6891248	       175 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GithubStatic            	 2676691	       442 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubStatic             	 2601064	       462 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubStatic      	  599408	      3151 ns/op	     296 B/op	       5 allocs/op
BenchmarkGoji_GithubStatic            	 1484044	       804 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_GithubStatic          	  138292	      8473 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGoRestful_GithubStatic       	   22000	     55267 ns/op	    4256 B/op	      13 allocs/op
BenchmarkGoJsonRest_GithubStatic      	  374398	      3920 ns/op	     329 B/op	      11 allocs/op
BenchmarkGorillaMux_GithubStatic      	   74012	     20612 ns/op	     976 B/op	       9 allocs/op
BenchmarkGowwwRouter_GithubStatic     	 3191044	       378 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GithubStatic      	 4895076	       243 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GithubStatic     	 4424451	       273 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_GithubStatic           	 8967478	       237 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_GithubStatic            	 4734001	       260 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubStatic         	  235126	      6771 ns/op	     736 B/op	       8 allocs/op
BenchmarkMartini_GithubStatic         	   49149	     26102 ns/op	     768 B/op	       9 allocs/op
BenchmarkPat_GithubStatic             	   25828	     39538 ns/op	    3648 B/op	      76 allocs/op
BenchmarkPossum_GithubStatic          	  600294	      2854 ns/op	     416 B/op	       3 allocs/op
BenchmarkProcyon_GithubStatic         	 5793034	       336 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GithubStatic        	 1000000	      1546 ns/op	     144 B/op	       4 allocs/op
BenchmarkRivet_GithubStatic           	 2371918	       503 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_GithubStatic           	  315771	      4810 ns/op	     248 B/op	       8 allocs/op
BenchmarkTigerTonic_GithubStatic      	 1000000	      1086 ns/op	      48 B/op	       1 allocs/op
BenchmarkTraffic_GithubStatic         	   25734	     45542 ns/op	    4664 B/op	      90 allocs/op
BenchmarkVulcan_GithubStatic          	  400029	      3422 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_GithubParam              	 1000000	      1375 ns/op	      96 B/op	       1 allocs/op
BenchmarkAero_GithubParam             	 2621391	       458 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubParam             	  413578	      4143 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_GithubParam            	  307480	      6932 ns/op	     544 B/op	       7 allocs/op
BenchmarkBone_GithubParam             	   41372	     26341 ns/op	    1888 B/op	      19 allocs/op
BenchmarkChi_GithubParam              	  461620	      4408 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_GithubParam  	 1640943	       731 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_GithubParam            	 1000000	      1371 ns/op	     128 B/op	       1 allocs/op
BenchmarkEcho_GithubParam             	 1624858	       785 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubParam              	 1408628	       859 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubParam       	  333366	      4935 ns/op	     712 B/op	       9 allocs/op
BenchmarkGoji_GithubParam             	  704654	      3387 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GithubParam           	  131739	      9809 ns/op	    1408 B/op	      13 allocs/op
BenchmarkGoJsonRest_GithubParam       	  230607	      6349 ns/op	     713 B/op	      14 allocs/op
BenchmarkGoRestful_GithubParam        	   17755	     69324 ns/op	    4352 B/op	      16 allocs/op
BenchmarkGorillaMux_GithubParam       	   42190	     32676 ns/op	    1296 B/op	      10 allocs/op
BenchmarkGowwwRouter_GithubParam      	  999433	      2627 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_GithubParam       	 1000000	      1127 ns/op	      96 B/op	       1 allocs/op
BenchmarkHttpTreeMux_GithubParam      	  545121	      2892 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_GithubParam            	  799413	      2345 ns/op	     128 B/op	       5 allocs/op
BenchmarkLARS_GithubParam             	 2040906	       587 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubParam          	  166533	     10029 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_GithubParam          	   35286	     38089 ns/op	    1152 B/op	      11 allocs/op
BenchmarkPat_GithubParam              	   35168	     31862 ns/op	    2408 B/op	      48 allocs/op
BenchmarkPossum_GithubParam           	  461620	      4195 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_GithubParam          	 2282331	       628 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GithubParam         	  704481	      2739 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_GithubParam            	 1000000	      1663 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_GithubParam            	  292683	      5601 ns/op	     344 B/op	       8 allocs/op
BenchmarkTigerTonic_GithubParam       	  127527	     11326 ns/op	    1176 B/op	      22 allocs/op
BenchmarkTraffic_GithubParam          	   44410	     40964 ns/op	    2816 B/op	      40 allocs/op
BenchmarkVulcan_GithubParam           	  307477	      4856 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_GithubAll                	    7053	    252400 ns/op	   13792 B/op	     167 allocs/op
BenchmarkAero_GithubAll               	   12816	     93645 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubAll               	    2031	    838160 ns/op	   86448 B/op	     943 allocs/op
BenchmarkBeego_GithubAll              	    1174	   1332978 ns/op	   88064 B/op	    1133 allocs/op
BenchmarkBone_GithubAll               	     100	  12026056 ns/op	  720160 B/op	    8620 allocs/op
BenchmarkChi_GithubAll                	    2725	    882884 ns/op	   87696 B/op	     609 allocs/op
BenchmarkCloudyKitRouter_GithubAll    	   10000	    108313 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_GithubAll              	    6302	    231485 ns/op	   20224 B/op	     167 allocs/op
BenchmarkEcho_GithubAll               	   10000	    155310 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubAll                	    6661	    179789 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubAll         	     957	   1171006 ns/op	  131656 B/op	    1686 allocs/op
BenchmarkGoji_GithubAll               	     768	   1743920 ns/op	   56112 B/op	     334 allocs/op
BenchmarkGojiv2_GithubAll             	     357	   3097242 ns/op	  352720 B/op	    4321 allocs/op
BenchmarkGoJsonRest_GithubAll         	     849	   1595919 ns/op	  134371 B/op	    2737 allocs/op
BenchmarkGoRestful_GithubAll          	     100	  12171531 ns/op	  910144 B/op	    2938 allocs/op
BenchmarkGorillaMux_GithubAll         	      69	  17884301 ns/op	  251651 B/op	    1994 allocs/op
BenchmarkGowwwRouter_GithubAll        	    3999	    574874 ns/op	   72144 B/op	     501 allocs/op
BenchmarkHttpRouter_GithubAll         	    9993	    190663 ns/op	   13792 B/op	     167 allocs/op
BenchmarkHttpTreeMux_GithubAll        	    2997	    599894 ns/op	   65856 B/op	     671 allocs/op
BenchmarkKocha_GithubAll              	    3240	    479554 ns/op	   23304 B/op	     843 allocs/op
BenchmarkLARS_GithubAll               	   10000	    111981 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubAll            	     788	   1535246 ns/op	  149409 B/op	    1624 allocs/op
BenchmarkMartini_GithubAll            	      79	  14894835 ns/op	  226555 B/op	    2325 allocs/op
BenchmarkPat_GithubAll                	      91	  16399848 ns/op	 1483152 B/op	   26963 allocs/op
BenchmarkPossum_GithubAll             	    2854	    733231 ns/op	   84448 B/op	     609 allocs/op
BenchmarkProcyon_GithubAll            	    9224	    156204 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GithubAll           	    2787	    618313 ns/op	   77328 B/op	     979 allocs/op
BenchmarkRivet_GithubAll              	    5702	    344998 ns/op	   16272 B/op	     167 allocs/op
BenchmarkTango_GithubAll              	    1377	   1167129 ns/op	   63825 B/op	    1618 allocs/op
BenchmarkTigerTonic_GithubAll         	     448	   2627153 ns/op	  193856 B/op	    4474 allocs/op
BenchmarkTraffic_GithubAll            	     100	  12629424 ns/op	  820744 B/op	   14114 allocs/op
BenchmarkVulcan_GithubAll             	    1873	    794544 ns/op	   19894 B/op	     609 allocs/op
```

## Gplus API Benchmarks
```
BenchmarkProcyon_GPlusStatic          	 4234502	       360 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlusParam           	 3701094	       411 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlusAll             	  214057	      7039 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlus2Params         	 2756644	       535 ns/op	       0 B/op	       0 allocs/op
```
```
BenchmarkAce_GPlusStatic              	 2782231	       453 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusStatic             	 6412588	       190 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusStatic             	 1000000	      1412 ns/op	     104 B/op	       3 allocs/op
BenchmarkBeego_GPlusStatic            	  330729	      4508 ns/op	     384 B/op	       7 allocs/op
BenchmarkBone_GPlusStatic             	 2493058	       687 ns/op	      32 B/op	       1 allocs/op
BenchmarkChi_GPlusStatic              	  749427	      2925 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_GPlusStatic  	 9033852	       166 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_GPlusStatic            	10704507	       111 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlusStatic             	 3809691	       324 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusStatic              	 3406670	       370 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusStatic       	  705384	      2700 ns/op	     280 B/op	       5 allocs/op
BenchmarkGoji_GPlusStatic             	 2439102	       614 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_GPlusStatic           	  157737	      8089 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGoJsonRest_GPlusStatic       	  545533	      3053 ns/op	     329 B/op	      11 allocs/op
BenchmarkGoRestful_GPlusStatic        	   91538	     19211 ns/op	    3872 B/op	      13 allocs/op
BenchmarkGorillaMux_GPlusStatic       	  181689	      7275 ns/op	     976 B/op	       9 allocs/op
BenchmarkGowwwRouter_GPlusStatic      	 8213703	       147 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GPlusStatic       	 9593307	       126 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GPlusStatic      	 6891672	       174 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_GPlusStatic            	 6876340	       175 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_GPlusStatic             	 4265600	       281 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusStatic          	  222066	      6023 ns/op	     736 B/op	       8 allocs/op
BenchmarkMartini_GPlusStatic          	  151789	     15232 ns/op	     768 B/op	       9 allocs/op
BenchmarkPat_GPlusStatic              	 1204353	      1070 ns/op	      96 B/op	       2 allocs/op
BenchmarkPossum_GPlusStatic           	  631345	      3104 ns/op	     416 B/op	       3 allocs/op
BenchmarkProcyon_GPlusStatic          	 4234502	       360 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GPlusStatic         	  749690	      1520 ns/op	     144 B/op	       4 allocs/op
BenchmarkRivet_GPlusStatic            	 5551627	       181 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_GPlusStatic            	  666195	      2497 ns/op	     200 B/op	       8 allocs/op
BenchmarkTigerTonic_GPlusStatic       	 3771216	       321 ns/op	      32 B/op	       1 allocs/op
BenchmarkTraffic_GPlusStatic          	  479652	      4570 ns/op	    1112 B/op	      16 allocs/op
BenchmarkVulcan_GPlusStatic           	 1277664	       916 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_GPlusParam               	 2426091	       478 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_GPlusParam              	11212353	       111 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusParam              	 1000000	      2349 ns/op	     480 B/op	       5 allocs/op
BenchmarkBeego_GPlusParam             	  303411	      6034 ns/op	     480 B/op	       7 allocs/op
BenchmarkBone_GPlusParam              	  428247	      4996 ns/op	     816 B/op	       6 allocs/op
BenchmarkChi_GPlusParam               	 1000000	      2354 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_GPlusParam   	 8327492	       144 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_GPlusParam             	 2540589	       429 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_GPlusParam              	 6213756	       175 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusParam               	 5629434	       212 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusParam        	 1000000	      1734 ns/op	     648 B/op	       8 allocs/op
BenchmarkGoji_GPlusParam              	 1411814	       849 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GPlusParam            	  413494	      3052 ns/op	    1328 B/op	      11 allocs/op
BenchmarkGoJsonRest_GPlusParam        	 1000000	      1896 ns/op	     649 B/op	      13 allocs/op
BenchmarkGoRestful_GPlusParam         	  148042	      7367 ns/op	    4192 B/op	      14 allocs/op
BenchmarkGorillaMux_GPlusParam        	  306972	      3909 ns/op	    1280 B/op	      10 allocs/op
BenchmarkGowwwRouter_GPlusParam       	 1386664	       873 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_GPlusParam        	 5560976	       650 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_GPlusParam       	  858264	      1958 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_GPlusParam             	 1000000	      1207 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_GPlusParam              	 4010539	       366 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusParam           	  164254	      9911 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_GPlusParam           	   61180	     19109 ns/op	    1072 B/op	      10 allocs/op
BenchmarkPat_GPlusParam               	  234994	      6052 ns/op	     576 B/op	      11 allocs/op
BenchmarkPossum_GPlusParam            	  571698	      4405 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_GPlusParam           	 3701094	       411 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GPlusParam          	  798232	      2716 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_GPlusParam             	 1333545	       881 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_GPlusParam             	  315570	      4591 ns/op	     264 B/op	       8 allocs/op
BenchmarkTigerTonic_GPlusParam        	  171318	      9843 ns/op	     856 B/op	      16 allocs/op
BenchmarkTraffic_GPlusParam           	   69729	     17018 ns/op	    1872 B/op	      21 allocs/op
BenchmarkVulcan_GPlusParam            	  461241	      3113 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_GPlus2Params             	 1000000	      1299 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_GPlus2Params            	 2700784	       440 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlus2Params            	  480110	      4015 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_GPlus2Params           	  249825	      7186 ns/op	     608 B/op	       7 allocs/op
BenchmarkBone_GPlus2Params            	   97491	     14146 ns/op	    1168 B/op	      10 allocs/op
BenchmarkChi_GPlus2Params             	  856506	      3285 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_GPlus2Params 	 3441732	       433 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_GPlus2Params           	 1000000	      1019 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_GPlus2Params            	 2224856	       606 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlus2Params             	 1463349	       812 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlus2Params      	  342668	      5020 ns/op	     712 B/op	       9 allocs/op
BenchmarkGoji_GPlus2Params            	  748707	      3459 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GPlus2Params          	   99910	     11521 ns/op	    1408 B/op	      14 allocs/op
BenchmarkGoJsonRest_GPlus2Params      	  266474	      7100 ns/op	     713 B/op	      14 allocs/op
BenchmarkGoRestful_GPlus2Params       	   46842	     24852 ns/op	    4384 B/op	      16 allocs/op
BenchmarkGorillaMux_GPlus2Params      	   52819	     26043 ns/op	    1296 B/op	      10 allocs/op
BenchmarkGowwwRouter_GPlus2Params     	  999399	      2905 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_GPlus2Params      	 1768501	       621 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_GPlus2Params     	  856536	      3075 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_GPlus2Params           	  705363	      2345 ns/op	     128 B/op	       5 allocs/op
BenchmarkLARS_GPlus2Params            	 2584138	       471 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlus2Params         	  147987	     10788 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_GPlus2Params         	   31639	     38745 ns/op	    1200 B/op	      13 allocs/op
BenchmarkPat_GPlus2Params             	   43759	     27128 ns/op	    2168 B/op	      33 allocs/op
BenchmarkPossum_GPlus2Params          	  461719	      4362 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_GPlus2Params         	 2756644	       535 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GPlus2Params        	  704233	      2791 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_GPlus2Params           	 1000000	      1376 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_GPlus2Params           	  333374	      5090 ns/op	     344 B/op	       8 allocs/op
BenchmarkTigerTonic_GPlus2Params      	   98407	     14277 ns/op	    1200 B/op	      22 allocs/op
BenchmarkTraffic_GPlus2Params         	   45771	     28815 ns/op	    2248 B/op	      28 allocs/op
BenchmarkVulcan_GPlus2Params          	  479324	      4039 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_GPlusAll                 	  134744	     11552 ns/op	     640 B/op	      11 allocs/op
BenchmarkAero_GPlusAll                	  443848	      4037 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusAll                	   26892	     43904 ns/op	    5488 B/op	      61 allocs/op
BenchmarkBeego_GPlusAll               	   15736	     79010 ns/op	    6016 B/op	      83 allocs/op
BenchmarkBone_GPlusAll                	   10000	    112152 ns/op	   11744 B/op	     109 allocs/op
BenchmarkChi_GPlusAll                 	   23242	     51156 ns/op	    5616 B/op	      39 allocs/op
BenchmarkCloudyKitRouter_GPlusAll     	  363162	      4069 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_GPlusAll               	  133246	     12013 ns/op	     672 B/op	      11 allocs/op
BenchmarkEcho_GPlusAll                	  168897	      7051 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusAll                 	  147992	      8135 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusAll          	   16963	     67429 ns/op	    8040 B/op	     103 allocs/op
BenchmarkGoji_GPlusAll                	   31065	     37943 ns/op	    3696 B/op	      22 allocs/op
BenchmarkGojiv2_GPlusAll              	   10000	    131317 ns/op	   17616 B/op	     154 allocs/op
BenchmarkGoJsonRest_GPlusAll          	   19524	     68784 ns/op	    8117 B/op	     170 allocs/op
BenchmarkGoRestful_GPlusAll           	   10000	    270213 ns/op	   55520 B/op	     192 allocs/op
BenchmarkGorillaMux_GPlusAll          	    8565	    208944 ns/op	   16112 B/op	     128 allocs/op
BenchmarkGowwwRouter_GPlusAll         	   37827	     34298 ns/op	    4752 B/op	      33 allocs/op
BenchmarkHttpRouter_GPlusAll          	  196658	      9362 ns/op	     640 B/op	      11 allocs/op
BenchmarkHttpTreeMux_GPlusAll         	   38186	     31881 ns/op	    4032 B/op	      38 allocs/op
BenchmarkKocha_GPlusAll               	   67368	     20008 ns/op	     976 B/op	      43 allocs/op
BenchmarkLARS_GPlusAll                	  235132	      5059 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusAll             	   13684	     91339 ns/op	    9568 B/op	     104 allocs/op
BenchmarkMartini_GPlusAll             	    5709	    297563 ns/op	   14016 B/op	     145 allocs/op
BenchmarkPat_GPlusAll                 	    9240	    164376 ns/op	   15264 B/op	     271 allocs/op
BenchmarkPossum_GPlusAll              	   24673	     42541 ns/op	    5408 B/op	      39 allocs/op
BenchmarkProcyon_GPlusAll             	  214057	      7039 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_GPlusAll            	   34455	     34109 ns/op	    5040 B/op	      63 allocs/op
BenchmarkRivet_GPlusAll               	  111050	     12674 ns/op	     768 B/op	      11 allocs/op
BenchmarkTango_GPlusAll               	   20101	     60247 ns/op	    3656 B/op	     104 allocs/op
BenchmarkTigerTonic_GPlusAll          	   10000	    123585 ns/op	   11600 B/op	     242 allocs/op
BenchmarkTraffic_GPlusAll             	   10000	    185116 ns/op	   26248 B/op	     341 allocs/op
BenchmarkVulcan_GPlusAll              	   40242	     44713 ns/op	    1274 B/op	      39 allocs/op
```

## Parse Benchmarks
```
BenchmarkProcyon_ParseStatic          	 4090732	       295 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_ParseParam           	11530256	       366 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_Parse2Params         	 4474801	       410 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_ParseAll             	  110059	     12338 ns/op	       0 B/op	       0 allocs/op
```
```
BenchmarkAce_ParseStatic              	 2584675	       469 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseStatic             	 5824732	       204 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseStatic             	  856536	      1563 ns/op	     120 B/op	       3 allocs/op
BenchmarkBeego_ParseStatic            	  399442	      4942 ns/op	     384 B/op	       7 allocs/op
BenchmarkBone_ParseStatic             	  630297	      2573 ns/op	     144 B/op	       3 allocs/op
BenchmarkChi_ParseStatic              	  922452	      2837 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_ParseStatic  	 6709532	       167 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_ParseStatic            	 9294118	       129 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_ParseStatic             	 3666870	       284 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseStatic              	 3155128	       381 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseStatic       	  598897	      3107 ns/op	     296 B/op	       5 allocs/op
BenchmarkGoji_ParseStatic             	 2001991	       759 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_ParseStatic           	  146197	      7768 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGoJsonRest_ParseStatic       	  545556	      3696 ns/op	     329 B/op	      11 allocs/op
BenchmarkGoRestful_ParseStatic        	   45076	     24553 ns/op	    4256 B/op	      13 allocs/op
BenchmarkGorillaMux_ParseStatic       	  139238	      9130 ns/op	     976 B/op	       9 allocs/op
BenchmarkGowwwRouter_ParseStatic      	 7888975	       148 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_ParseStatic       	 9517041	       126 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_ParseStatic      	 4523989	       265 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_ParseStatic            	 6212488	       193 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_ParseStatic             	 4135358	       291 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseStatic          	  260684	      6044 ns/op	     736 B/op	       8 allocs/op
BenchmarkMartini_ParseStatic          	   86894	     15893 ns/op	     768 B/op	       9 allocs/op
BenchmarkPat_ParseStatic              	  499645	      2577 ns/op	     240 B/op	       5 allocs/op
BenchmarkPossum_ParseStatic           	 1000000	      2633 ns/op	     416 B/op	       3 allocs/op
BenchmarkProcyon_ParseStatic          	 4090732	       295 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_ParseStatic         	 1000000	      1469 ns/op	     144 B/op	       4 allocs/op
BenchmarkRivet_ParseStatic            	 3868250	       308 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_ParseStatic            	  399403	      4035 ns/op	     248 B/op	       8 allocs/op
BenchmarkTigerTonic_ParseStatic       	 1000000	      1000 ns/op	      48 B/op	       1 allocs/op
BenchmarkTraffic_ParseStatic          	  146238	     10114 ns/op	    1256 B/op	      19 allocs/op
BenchmarkVulcan_ParseStatic           	  600294	      2348 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_ParseParam               	 1000000	      1151 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_ParseParam              	 4629920	       259 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseParam              	  571825	      3071 ns/op	     467 B/op	       5 allocs/op
BenchmarkBeego_ParseParam             	  272530	      6083 ns/op	     416 B/op	       7 allocs/op
BenchmarkBone_ParseParam              	  324084	      5793 ns/op	     896 B/op	       7 allocs/op
BenchmarkChi_ParseParam               	  705363	      3810 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_ParseParam   	 5765114	       208 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_ParseParam             	 1505602	       801 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_ParseParam              	 3031015	       395 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseParam               	 2834872	       427 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseParam        	  428274	      4892 ns/op	     664 B/op	       8 allocs/op
BenchmarkGoji_ParseParam              	  799407	      2949 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_ParseParam            	  130344	      9167 ns/op	    1360 B/op	      12 allocs/op
BenchmarkGoJsonRest_ParseParam        	  428268	      5721 ns/op	     649 B/op	      13 allocs/op
BenchmarkGoRestful_ParseParam         	   33216	     40998 ns/op	    4576 B/op	      14 allocs/op
BenchmarkGorillaMux_ParseParam        	  143649	      7893 ns/op	    1280 B/op	      10 allocs/op
BenchmarkGowwwRouter_ParseParam       	 1000000	      1698 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_ParseParam        	 3506293	       344 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_ParseParam       	 1000000	      1088 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_ParseParam             	 2456520	       499 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_ParseParam              	 7837394	       135 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseParam           	  499672	      4220 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_ParseParam           	  170085	      7218 ns/op	    1072 B/op	      10 allocs/op
BenchmarkPat_ParseParam               	  444123	      3689 ns/op	     992 B/op	      15 allocs/op
BenchmarkPossum_ParseParam            	  999466	      1443 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_ParseParam           	11530256	       366 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_ParseParam          	  922452	      1999 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_ParseParam             	 1398692	       738 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_ParseParam             	  500163	      3309 ns/op	     280 B/op	       8 allocs/op
BenchmarkTigerTonic_ParseParam        	  333093	      7723 ns/op	     784 B/op	      15 allocs/op
BenchmarkTraffic_ParseParam           	   79930	     13849 ns/op	    1896 B/op	      21 allocs/op
BenchmarkVulcan_ParseParam            	  631146	      1906 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_Parse2Params             	 1707183	       900 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_Parse2Params            	 4392962	       254 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Parse2Params            	  499635	      3366 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_Parse2Params           	  444052	      6419 ns/op	     480 B/op	       7 allocs/op
BenchmarkBone_Parse2Params            	  352998	      5577 ns/op	     848 B/op	       6 allocs/op
BenchmarkChi_Parse2Params             	 1000000	      2454 ns/op	     432 B/op	       3 allocs/op
BenchmarkCloudyKitRouter_Parse2Params 	 7156608	       164 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_Parse2Params           	 2635766	       448 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_Parse2Params            	 5709801	       202 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Parse2Params             	 5764784	       207 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Parse2Params      	  922423	      2010 ns/op	     712 B/op	       9 allocs/op
BenchmarkGoji_Parse2Params            	 1304047	       918 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Parse2Params          	  399500	      3079 ns/op	    1344 B/op	      11 allocs/op
BenchmarkGoJsonRest_Parse2Params      	  855309	      2159 ns/op	     713 B/op	      14 allocs/op
BenchmarkGoRestful_Parse2Params       	  124894	     27347 ns/op	    4928 B/op	      14 allocs/op
BenchmarkGorillaMux_Parse2Params      	   96703	     12729 ns/op	    1296 B/op	      10 allocs/op
BenchmarkGowwwRouter_Parse2Params     	 1000000	      2627 ns/op	     432 B/op	       3 allocs/op
BenchmarkHttpRouter_Parse2Params      	 1620472	       755 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Parse2Params     	  705330	      2934 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_Parse2Params           	  521845	      2250 ns/op	     128 B/op	       5 allocs/op
BenchmarkLARS_Parse2Params            	 3106609	       379 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Parse2Params         	  168823	     10217 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_Parse2Params         	   62364	     19746 ns/op	    1152 B/op	      11 allocs/op
BenchmarkPat_Parse2Params             	  230578	      9737 ns/op	     752 B/op	      16 allocs/op
BenchmarkPossum_Parse2Params          	  363289	      4008 ns/op	     496 B/op	       5 allocs/op
BenchmarkProcyon_Parse2Params         	 4474801	       410 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_Parse2Params        	  750835	      2966 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Parse2Params           	 1000000	      1229 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_Parse2Params           	  374733	      4824 ns/op	     312 B/op	       8 allocs/op
BenchmarkTigerTonic_Parse2Params      	  119913	     13889 ns/op	    1168 B/op	      22 allocs/op
BenchmarkTraffic_Parse2Params         	   69717	     15644 ns/op	    1944 B/op	      22 allocs/op
BenchmarkVulcan_Parse2Params          	  705384	      2903 ns/op	      98 B/op	       3 allocs/op

BenchmarkAce_ParseAll                 	   64159	     22182 ns/op	     640 B/op	      16 allocs/op
BenchmarkAero_ParseAll                	  159886	      7580 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseAll                	   13872	     80221 ns/op	    8928 B/op	     110 allocs/op
BenchmarkBeego_ParseAll               	   10000	    133228 ns/op	    9952 B/op	     114 allocs/op
BenchmarkBone_ParseAll                	   10000	    134185 ns/op	   16208 B/op	     147 allocs/op
BenchmarkChi_ParseAll                 	   16479	    104718 ns/op	   11232 B/op	      78 allocs/op
BenchmarkCloudyKitRouter_ParseAll     	  235138	      5919 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_ParseAll               	   74020	     16486 ns/op	     928 B/op	      16 allocs/op
BenchmarkEcho_ParseAll                	  162093	      7343 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseAll                 	   81012	     14811 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseAll          	   10000	    123252 ns/op	   13728 B/op	     181 allocs/op
BenchmarkGoji_ParseAll                	   18294	     58816 ns/op	    5376 B/op	      32 allocs/op
BenchmarkGojiv2_ParseAll              	   10000	    212961 ns/op	   34448 B/op	     277 allocs/op
BenchmarkGoJsonRest_ParseAll          	   10000	    127673 ns/op	   13866 B/op	     321 allocs/op
BenchmarkGoRestful_ParseAll           	    1790	    713113 ns/op	  117600 B/op	     354 allocs/op
BenchmarkGorillaMux_ParseAll          	    5450	    414860 ns/op	   30288 B/op	     250 allocs/op
BenchmarkGowwwRouter_ParseAll         	   26145	     52848 ns/op	    6912 B/op	      48 allocs/op
BenchmarkHttpRouter_ParseAll          	  144478	     12110 ns/op	     640 B/op	      16 allocs/op
BenchmarkHttpTreeMux_ParseAll         	   25789	     44104 ns/op	    5728 B/op	      51 allocs/op
BenchmarkKocha_ParseAll               	   48939	     25671 ns/op	    1112 B/op	      54 allocs/op
BenchmarkLARS_ParseAll                	  131816	      8973 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseAll             	    9992	    178633 ns/op	   19136 B/op	     208 allocs/op
BenchmarkMartini_ParseAll             	    3525	    486300 ns/op	   25072 B/op	     253 allocs/op
BenchmarkPat_ParseAll                 	   10000	    157719 ns/op	   15216 B/op	     308 allocs/op
BenchmarkPossum_ParseAll              	   13879	     89427 ns/op	   10816 B/op	      78 allocs/op
BenchmarkProcyon_ParseAll             	  110059	     12338 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_ParseAll            	   18242	     68196 ns/op	    8352 B/op	     120 allocs/op
BenchmarkRivet_ParseAll               	   56570	     21012 ns/op	     912 B/op	      16 allocs/op
BenchmarkTango_ParseAll               	   10000	    121014 ns/op	    7168 B/op	     208 allocs/op
BenchmarkTigerTonic_ParseAll          	    8564	    182071 ns/op	   16048 B/op	     332 allocs/op
BenchmarkTraffic_ParseAll             	    4611	    406753 ns/op	   45520 B/op	     605 allocs/op
BenchmarkVulcan_ParseAll              	   14206	     76412 ns/op	    2548 B/op	      78 allocs/op
```

## Static Benchmarks
```
BenchmarkProcyon_StaticAll            	   15219	     78365 ns/op	       0 B/op	       0 allocs/op
```
```
BenchmarkAce_StaticAll                	   10000	    114881 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_StaticAll               	   26060	     46055 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpServeMux_StaticAll       	    8565	    139838 ns/op	       0 B/op	       0 allocs/op
BenchmarkBeego_StaticAll              	    1498	    845290 ns/op	   71008 B/op	    1097 allocs/op
BenchmarkBear_StaticAll               	    4441	    344295 ns/op	   20272 B/op	     469 allocs/op
BenchmarkBone_StaticAll               	    3996	    290996 ns/op	       0 B/op	       0 allocs/op
BenchmarkChi_StaticAll                	    3426	    565135 ns/op	   67824 B/op	     471 allocs/op
BenchmarkCloudyKitRouter_StaticAll    	   17324	     72966 ns/op	       0 B/op	       0 allocs/op
BenchmarkDenco_StaticAll              	   39315	     30926 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_StaticAll               	   10000	    107349 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_StaticAll                	   10000	    113708 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_StaticAll         	    3630	    491297 ns/op	   46312 B/op	     785 allocs/op
BenchmarkGoji_StaticAll               	    6310	    185160 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_StaticAll             	     772	   1581370 ns/op	  205984 B/op	    1570 allocs/op
BenchmarkGoJsonRest_StaticAll         	    1965	    797685 ns/op	   51653 B/op	    1727 allocs/op
BenchmarkGoRestful_StaticAll          	     190	   6393030 ns/op	  613280 B/op	    2053 allocs/op
BenchmarkGorillaMux_StaticAll         	     338	   4085444 ns/op	  153233 B/op	    1413 allocs/op
BenchmarkGowwwRouter_StaticAll        	   13783	     87076 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_StaticAll         	   23158	     51378 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_StaticAll        	   23283	     47079 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_StaticAll              	   24373	     60592 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_StaticAll               	   16095	     76247 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_StaticAll            	    1816	    929257 ns/op	  115552 B/op	    1256 allocs/op
BenchmarkMartini_StaticAll            	     218	   6797197 ns/op	  125446 B/op	    1717 allocs/op
BenchmarkPat_StaticAll                	     182	   8466688 ns/op	  602832 B/op	   12559 allocs/op
BenchmarkPossum_StaticAll             	    4616	    516543 ns/op	   65312 B/op	     471 allocs/op
BenchmarkProcyon_StaticAll            	   15219	     78365 ns/op	       0 B/op	       0 allocs/op
BenchmarkR2router_StaticAll           	    5718	    247754 ns/op	   22608 B/op	     628 allocs/op
BenchmarkRivet_StaticAll              	    9223	    126558 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_StaticAll              	    1736	    777636 ns/op	   39209 B/op	    1256 allocs/op
BenchmarkTigerTonic_StaticAll         	   10000	    190739 ns/op	    7376 B/op	     157 allocs/op
BenchmarkTraffic_StaticAll            	     140	   9035511 ns/op	  754863 B/op	   14601 allocs/op
BenchmarkVulcan_StaticAll             	    2787	    573482 ns/op	   15386 B/op	     471 allocs/op
```

## Realistic Benchmarks 
As you know, We need the recovery to prevent the app crashing and an unique id like uuid for logging while developing our applications. 
What Procyon makes difference is to support these features without using something like middleware. As other frameworks don't support 
directly and you have to use middleware, it makes them slower. Here are the benchmark results.

## Benchmarks with Recovery
```
BenchmarkAero_Param                   	 2903332	       410 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_Param                	 2343621	       499 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param                    	 1768957	       683 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_Param                   	 1000000	      1134 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_Param5                  	 2256878	       515 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_Param5               	 2128104	       558 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param5                   	 1206388	       989 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_Param5                  	  924157	      1615 ns/op	      48 B/op	       1 allocs/op

BenchmarkProcyon_Param20              	 1397652	       849 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param20                 	  922636	      1126 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param20                  	  544518	      2247 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_Param20                 	  479668	      3059 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_GithubStatic            	 2990430	       398 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GithubStatic         	 2275324	       521 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubStatic             	 1637776	       734 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GithubStatic            	 1000000	      1225 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_GithubParam             	 1920070	       624 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GithubParam          	 1422472	       864 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubParam              	 1000000	      1131 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GithubParam             	 1000000	      1627 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_GithubAll               	    9021	    132993 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GithubAll            	    5708	    201223 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubAll                	    5995	    243875 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GithubAll               	    4441	    336457 ns/op	    9744 B/op	     203 allocs/op

BenchmarkAero_GPlusStatic             	 3248190	       363 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlusStatic          	 2573286	       468 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusStatic              	 1873660	       637 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlusStatic             	 1000000	      1058 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_GPlusParam              	 2584456	       469 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlusParam           	 1998664	       598 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusParam               	 1473078	       811 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlusParam              	 1000000	      1283 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_GPlus2Params            	 1956105	       613 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlus2Params         	 1589557	       746 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlus2Params             	 1000000	      1049 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlus2Params            	 1000000	      1394 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_GPlusAll                	  171237	      7000 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_GPlusAll             	  117562	      9877 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusAll                 	   97491	     12411 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlusAll                	   82716	     17122 ns/op	     624 B/op	      13 allocs/op

BenchmarkAero_ParseStatic             	 3138910	       382 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_ParseStatic          	 1909386	       628 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseStatic              	 1814361	       655 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_ParseStatic             	 1000000	      1045 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_ParseParam              	 2412782	       503 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_ParseParam           	 2028723	       585 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseParam               	 1654003	       739 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_ParseParam              	 1000000	      1190 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_Parse2Params            	 2310505	       572 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_Parse2Params         	 1853401	       658 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Parse2Params             	 1362416	       869 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_Parse2Params            	  922452	      1318 ns/op	      48 B/op	       1 allocs/op

BenchmarkAero_ParseAll                	   94404	     12711 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_ParseAll             	   71793	     19905 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseAll                 	   53054	     23033 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_ParseAll                	   40783	     32129 ns/op	    1248 B/op	      26 allocs/op

BenchmarkAero_StaticAll               	   15961	     75128 ns/op	       0 B/op	       0 allocs/op
BenchmarkProcyon_StaticAll            	   10000	    115427 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_StaticAll                	    7494	    161097 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_StaticAll               	    6662	    231263 ns/op	    7536 B/op	     157 allocs/op
```

## Benchmarks with Recovery-UUID
```
BenchmarkProcyon_Param                	 1700942	       680 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param                   	  790000	      1624 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_Param                    	  858344	      1651 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_Param                   	  706218	      2064 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_Param5               	 1653933	       723 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param5                  	  858374	      1505 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_Param5                   	  705400	      2057 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_Param5                  	  570409	      2830 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_Param20              	 1000000	      1017 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param20                 	  676702	      2022 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_Param20                  	  374724	      3381 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_Param20                 	  342613	      4416 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_GithubStatic         	 2167550	       704 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubStatic            	 1000000	      1262 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_GithubStatic             	  797956	      1626 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_GithubStatic            	  665088	      2302 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_GithubParam          	 1458230	      1040 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubParam             	 1000000	      1468 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_GithubParam              	  599768	      2314 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_GithubParam             	  444608	      2749 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_GithubAll            	    4790	    245404 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubAll               	    4134	    315785 ns/op	    9744 B/op	     203 allocs/op
BenchmarkGin_GithubAll                	    2726	    440323 ns/op	    9744 B/op	     203 allocs/op
BenchmarkEcho_GithubAll               	    2606	    583256 ns/op	   22736 B/op	     609 allocs/op

BenchmarkProcyon_GPlusStatic          	 1835948	       652 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusStatic             	 1000000	      1293 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_GPlusStatic              	  922423	      1582 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_GPlusStatic             	  704410	      2224 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_GPlusParam           	 1545254	       776 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusParam              	 1000000	      1447 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_GPlusParam               	  749469	      1872 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_GPlusParam              	  631130	      2425 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_GPlus2Params         	 1345892	       897 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlus2Params            	 1000000	      1596 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_GPlus2Params             	  856456	      2101 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_GPlus2Params            	  571011	      2475 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_GPlusAll             	   98289	     13062 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusAll                	   72674	     17942 ns/op	     624 B/op	      13 allocs/op
BenchmarkGin_GPlusAll                 	   50805	     27348 ns/op	     624 B/op	      13 allocs/op
BenchmarkEcho_GPlusAll                	   33907	     35159 ns/op	    1456 B/op	      39 allocs/op

BenchmarkProcyon_ParseStatic          	 1841326	       653 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseStatic             	  922486	      1265 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_ParseStatic              	 1000000	      1638 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_ParseStatic             	  705400	      2024 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_ParseParam           	 1589708	       752 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseParam              	 1000000	      1325 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_ParseParam               	  748152	      1808 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_ParseParam              	  666133	      2342 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_Parse2Params         	 1499960	       796 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Parse2Params            	 1000000	      1433 ns/op	      48 B/op	       1 allocs/op
BenchmarkGin_Parse2Params             	  749568	      1896 ns/op	      48 B/op	       1 allocs/op
BenchmarkEcho_Parse2Params            	  631176	      2568 ns/op	     112 B/op	       3 allocs/op

BenchmarkProcyon_ParseAll             	   50389	     24686 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseAll                	   36228	     36975 ns/op	    1248 B/op	      26 allocs/op
BenchmarkGin_ParseAll                 	   27190	     51420 ns/op	    1248 B/op	      26 allocs/op
BenchmarkEcho_ParseAll                	   19657	     62686 ns/op	    2912 B/op	      78 allocs/op

BenchmarkProcyon_StaticAll            	   10000	    140768 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_StaticAll               	    5994	    222543 ns/op	    7536 B/op	     157 allocs/op
BenchmarkGin_StaticAll                	    4285	    338801 ns/op	    7536 B/op	     157 allocs/op
BenchmarkEcho_StaticAll               	    3632	    441390 ns/op	   17584 B/op	     471 allocs/op
```

# Features
Procyon offers several features to make development process easier. 
Thanks to them, you code your application and deploy fast. The most important feature 
offered by Procyon is **Dependency Injection**. While application starts, instances are automatically
injected.

Here are all list of features offered:
* Dependency Injection
* Project Structure
* Routing
* Logger
* Configurable
* Fast
* Events
* Interceptors
* Error Handling
* Request and Path Variable Binding
* Extendable

## Dependency Injection
Dependency injection is not supported by standard libraries in Go. However, A few libraries like [Wire](https://github.com/google/wire) 
developed by Google support it at compile-time while Procyon offers it at runtime. It's the first http 
framework supporting the dependency injection at runtime. Instances created by procyon are called '**Peas**'.

Note that you need to register construction functions by using the function **core.Register**
because Go doesn't support something like annotation in Java.

The example is given below.

```go
package main

type Struct1 struct {
    // fields...
}

func newStruct1() Struct1 {
    return Struct1{}
}

type Struct2 struct {
    struct1 Struct1
    // fields ...
}

func newStruct2(struct1 Struct1) Struct2 {
    return Struct2{
        struct1,
    }
}

func init() {
	core.register(newStruct1)
    core.register(newStruct2)
}

...

```

## Project Structure
Thanks to the dependency injection, you can organize your project structure like Controller-Service-Repository.
You can look into [the test application](https://github.com/procyon-projects/procyon-test-app) to understand
how to use.

## Routing
You can easily register your handler methods by using **HandlerRegistry** in the method **RegisterHandlers**.

```go
func (controller HelloWorldController) RegisterHandlers(registry web.HandlerRegistry) {
	registry.Register(
		web.NewHandler(controller.HelloWorld, web.WithPath("/api/helloworld"))
	)
}
```


## Logger
Procyon offers a standard logger for logging. It can be got through dependency injection. 

The example is given below.

```go
package main

import (
	context "github.com/procyon-projects/procyon-context"
)

type TestController struct {
    logger     context.Logger
}

func newTestController(logger context.Logger) TestController {
    return TestController{}
}
```

## Configurable
Procyon application can be easily configured because it is designed according to configuration
properties. For example, when you want application to start on port 3000, what you need to do is
to specify the command-line parameter **--server.port** as 3030.  Procyon application contains many
configuration properties like this.


Also, You can define your configuration struct and its fields are automatically filled.

Note that you need to register your configuration struct.

The example is given below.
```go
package main

type MyProperties struct {
	Port              int       `yaml:"port" json:"port" default:"8080"`
    ApplicationName   string    `yaml:"name" json:"name" default:"Test Application"`
}

func newMyProperties() *MyProperties {
	return &MyProperties{}
}

func (properties *MyProperties) GetPrefix() string {
	return "myproperties"
}

func init() {
	core.register(newMyProperties)
}

```

## Fast
You can look into [the benchmark result](benchmarks.md) to find out how fast Procyon is.

## Events
The integration of events hasn't been completed yet.

## Interceptors
The integration of interceptors hasn't been completed yet.

## Error Handling
The integration of error handling hasn't been completed yet.

## Request and Path Variable Binding
The integration of error handling hasn't been completed yet.

## Extendable
It's easier to extend the application. You can look into the [modules](modules.md) to find out
how to do.

## How to contribute to Procyon?
* Contribute to our projects and become a member of our team
* Report bugs you find

## License
Procyon Framework is released under version 2.0 of the Apache License
