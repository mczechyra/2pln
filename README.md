# 2pln
------

Simple command line tool to convert a given amount of currency according to the current exchange rate from api.exchangeratesapi.io  
After first download data current status is save in local temp folder, so there is no need to download data again.  
Cache is refreshed after every day change - only one time after first request.

**Compile:**  
No external dependencies, project uses pure golang stamdard library  
To compile type:  
go build

**Usage:**  
$ 2pln 99 USD  
[2020-09-18] 99.0000 USD = 373.1596 PLN  

**Output meaning:**  
[2020-09-18] - date of currency exchange rate  
99.0000 USD - input amount of currency  
373.1596 PLN - calculated result in PLN