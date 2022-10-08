SGLL - Simple Go Logging Library

Easy to use thread safe go logging library.
Features:
- 4 log levels (Debug, Info, Warn, Error)
- variables support
- logging to file

Example
`
    import logging "github.com/paweljvm/sgl"

    func main() {
        // start logging to file
        // if you don't call that you will log only to the console
        logging.LogToFile("out.log")
        // remember to close logging file when you're done
        defer logging.CloseLogFile()
		logging.SetLevel(logging.DEBUG)

        logging.Info("This info example message without variables")
        logging.Warn("This is {} example message without variables", "warn")        
        logging.Error("This is error: {} ", error.Error())
        logging.Error("Debug {} {} {}", 1, 2, 3)      
    }
`