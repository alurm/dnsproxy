package main

import ("encoding/json"; "log"; "fmt"; "os"; "io")

type configuration struct { Server string; Blacklist []string }

func main() {
	var c configuration
	{
		var bytes []byte
		{
			file, err := os.Open("configuration.json")
			if err != nil { log.Fatal(err) }
			defer file.Close()
			bytes, err = io.ReadAll(file)
		}

		err := json.Unmarshal(
			bytes,
			&c,
		)
		if err != nil { log.Fatal(err) }
	}
	fmt.Printf("%#v\n", c)
}
