package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

func main() {

}

/*
	This function outlines the improper way to use waitgroups and channels
	Client - http client can be mocked
	timeout - time in seconds passes as an integer
*/
func buggy_function(client *http.Client, timeout time.Duration) {

	service_endpoints := []string{"http://someservice.com", "http://someservice.com", "http://someservice.com", "http://someservice.com", "http://someservice.com"}

	wg := sync.WaitGroup{}
	wg.Add(len(service_endpoints))

	type response_data struct {
		response    string
		err         error
		status_code int
	}
	response_chan := make(chan *response_data, len(service_endpoints)) // buffered channel
	defer close(response_chan)                                         // this is an issue, as the function can return before all the goroutines respond
	for _, endpoint := range service_endpoints {
		go func(wg sync.WaitGroup, response_chan chan *response_data, timeout time.Duration) {
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				wg.Done()
				response_chan <- &response_data{
					err:         err,
					status_code: 500,
				}
				return
			}
			ctx, cancelFunc := context.WithTimeout(req.Context(), timeout)
			defer cancelFunc()
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				wg.Done()
				response_chan <- &response_data{
					err:         err,
					status_code: 500,
				}
				return
			}
			if resp.StatusCode >= http.StatusBadRequest {
				wg.Done()
				response_chan <- &response_data{
					err:         errors.New("non 2xx response recieved from server"),
					status_code: resp.StatusCode,
				}
				return
			}

			response_chan <- &response_data{
				response:    "success",
				status_code: resp.StatusCode,
			}
			return
		}(wg, response_chan, timeout)
	}

	wg.Wait()
	var recievedResponses []*response_data
	for i := 0; i < len(service_endpoints); i++ {
		select {
		case response := <-response_chan:
			recievedResponses = append(recievedResponses, response)
			if response.err != nil {

			}
		default:
			if len(recievedResponses) == len(service_endpoints) {
				break
			}
		}
	}

}

func fixed(client *http.Client, timeout time.Duration) {
	service_endpoints := []string{"http://someservice.com", "http://someservice.com", "http://someservice.com", "http://someservice.com", "http://someservice.com"}

	wg := sync.WaitGroup{}
	wg.Add(len(service_endpoints))

	type response_data struct {
		response    string
		err         error
		status_code int
	}
	response_chan := make(chan *response_data, len(service_endpoints))
	for _, endpoint := range service_endpoints {
		go func(wg sync.WaitGroup, response_chan chan *response_data, timeout time.Duration) {
			defer wg.Done() // use defer here
			defer func() {
				if err := recover(); err != nil {
					// log out or do something here to indicate the program panics
					return
				}
			}()
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				response_chan <- &response_data{
					err:         err,
					status_code: 500,
				}

				return
			}
			ctx, cancelFunc := context.WithTimeout(req.Context(), timeout)
			defer cancelFunc()
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				response_chan <- &response_data{
					err:         err,
					status_code: 500,
				}
				return
			}
			if resp.StatusCode >= http.StatusBadRequest {
				response_chan <- &response_data{
					err:         errors.New("non 2xx response recieved from server"),
					status_code: resp.StatusCode,
				}
				return
			}
			response_chan <- &response_data{
				response:    "success",
				status_code: resp.StatusCode,
			}
			return
		}(wg, response_chan, timeout)
	}

	go func(wg sync.WaitGroup) { // moved this to it's own go routine so that we can go into the infinite for loop
		wg.Wait()
		close(response_chan)
	}(wg)

	var recievedResponses []*response_data
	for {
		select { // in this case, we only have one case to look at, when we get a response, and we break out when we recieved all the responses as we expected. ( the buffered channel will panic if we try to write more than the buffer )
		case response := <-response_chan:
			recievedResponses = append(recievedResponses, response)
			if response.err != nil {

			}
			if len(recievedResponses) == len(service_endpoints) {
				break
			}
		}
		break
	}
}

func fixedk(client *http.Client, timeout time.Duration){
	service_endpoints := &#91;]string{"http://someservice.com","http://someservice.com","http://someservice.com","http://someservice.com","http://someservice.com"}

	wg := sync.WaitGroup{}
	wg.Add(len(service_endpoints))
	
	type response_data struct{
		response string 
		err error
		status_code int
	}
	response_chan := make(chan *response_data, len(service_endpoints))
	for _, endpoint := range service_endpoints{
		go func (wg sync.WaitGroup, response_chan chan *response_data, timeout time.Duration){

			defer func ()  {
				if err := recover(); err != nil {
					// log out or do something here to indicate the program panics
					return 
				}
			}()
                        defer wg.Done() // use defer here 
			req, err  := http.NewRequest(http.MethodGet,endpoint,nil)
			if err != nil {
				response_chan &lt;- &amp;response_data{
					err : err ,
					status_code: 500,
				}
			
				return
			}
			ctx, cancelFunc := context.WithTimeout(req.Context(),timeout)
			defer cancelFunc()
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				response_chan &lt;- &amp;response_data{
					err : err ,
					status_code: 500,
				}
				return 
			}
			if resp.StatusCode >= http.StatusBadRequest{
				response_chan &lt;- &amp;response_data{
					err : errors.New("non 2xx response recieved from server") ,
					status_code: resp.StatusCode,
				}
				return
			}
			response_chan &lt;- &amp;response_data{
				response: "success",
				status_code: resp.StatusCode,
			}
			return
		}(wg, response_chan,timeout)
	}
	
	go func (wg sync.WaitGroup)  {  // moved this to it's own go routine so that we can go into the infinite for loop 
		wg.Wait()
		close(response_chan)
	}(wg)

	var recievedResponses &#91;]*response_data
	for {
		select{ // in this case, we only have one case to look at, when we get a response, and we break out when we recieved all the responses as we expected. ( the buffered channel will panic if we try to write more than the buffer )
		case response := &lt;- response_chan: 
		recievedResponses = append(recievedResponses, response)
			if response.err != nil {
	
			}
			if len(recievedResponses) == len(service_endpoints){
				break
			}
		}
		break
	}
}