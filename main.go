package main

import (
    "reflect"
    "bufio"
	"fmt"
	"strings"
	"strconv"
	"compress/gzip"
    "log"
	"os"
	"path/filepath"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/protobuf/proto"

	proto_struct "github.com/balabanas/memc-protobuf/proto"
)

func insertAppsinstalled(memcAddr string, key string, ua *proto_struct.UserApps, dryRun bool) (bool) {
    packed, err := proto.Marshal(ua)
        if err != nil {
            panic(err)
    }
    if dryRun {
        log.Printf("%s - %s -> %s", memcAddr, key, ua)
        return true
    }
    memc := memcache.New(memcAddr)
    if err := memc.Set(&memcache.Item{Key: key, Value: packed, Flags: 0, Expiration: 0}); err != nil {
        log.Printf("Cannot write to memc %s: %v", memcAddr, err)
        return false
    }
    return true
}


func parse(memcAddr string, options *Options, dryRun bool) (bool) {
    files, err := filepath.Glob(options.pattern)
	    if err != nil {
		    log.Fatal(err)
	    }
	for _, fn := range files {
		processed := 0
		errors := 0
		log.Printf("Processing %s %s %s", fn, processed, errors)
		fd, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		gzReader, err := gzip.NewReader(fd)
		if err != nil {
			log.Fatal(err)
		}
		defer gzReader.Close()

        scanner := bufio.NewScanner(gzReader)

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, "\t")
			if len(parts) != 5 {
				errors++
				log.Printf("Invalid line format: %s", line)
				continue
			}

			devType := parts[0]
			devID := parts[1]
			lat, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				errors++
				log.Printf("Error parsing latitude: %s", line)
				continue
			}
			lon, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				errors++
				log.Printf("Error parsing longitude: %s", line)
				continue
			}

			rawApps := strings.Split(parts[4], ",")
			var apps []uint32
			for _, rawApp := range rawApps {
				app, err := strconv.Atoi(rawApp)
				if err == nil {
					apps = append(apps, uint32(app))
				}
			}
            ua := &proto_struct.UserApps{
                Lat:  proto.Float64(lat),
                Lon:  proto.Float64(lon),
                Apps: apps,
            }
            key := fmt.Sprintf("%s:%s", devType, devID)
            result := insertAppsinstalled(memcAddr, key, ua, false)
            fmt.Println(result)

            processed++
		}

        if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Processed %d lines in %s\n", processed, fn)
		fmt.Printf("Encountered %d errors in %s\n", errors, fn)

	}
    return true

}


func prototest() {
    sample := "idfa\t1rfw452y52g2gq4g\t55.55\t42.42\t1423,43,567,3,7,23\ngaid\t7rfw452y52g2gq4g\t55.55\t42.42\t7423,424"
	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")

		//devType := parts[0]
		//devID := parts[1]
		lat, _ := strconv.ParseFloat(parts[2], 64)
		lon, _ := strconv.ParseFloat(parts[3], 64)
		rawApps := strings.Split(parts[4], ",")

		var apps []uint32
		for _, rawApp := range rawApps {
			if app, err := strconv.Atoi(rawApp); err == nil {
				apps = append(apps, uint32(app))
			}
		}

		ua := &proto_struct.UserApps{
			Lat:  proto.Float64(lat),
			Lon:  proto.Float64(lon),
			Apps: apps,
		}

		packed, err := proto.Marshal(ua)
		if err != nil {
			panic(err)
		}

		unpacked := &proto_struct.UserApps{}
		if err := proto.Unmarshal(packed, unpacked); err != nil {
			panic(err)
		}

		if proto.Equal(ua, unpacked) {
			fmt.Println("Protobuf objects are equal.")
		} else {
			fmt.Println("Protobuf objects are NOT equal.")
		}


	}

}


func main() {
   	options := &Options{
		idfa:  "your_idfa_option_value",
		gaid:  "your_gaid_option_value",
		adid:  "your_adid_option_value",
		dvid:  "your_dvid_option_value",
		pattern: "data/*.tsv.gz", // Replace options.pattern with the desired pattern
	}
// 	deviceMemc := map[string]string{
// 		"idfa": options.idfa,
// 		"gaid": options.gaid,
// 		"adid": options.adid,
// 		"dvid": options.dvid,
// 	}


     mc := memcache.New("localhost:33013", "localhost:33014", "localhost:33015", "localhost:33016")
// 	 mcd := memcache.New("localhost:33015", "localhost:33016")
//      mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
// 	 mcd.Set(&memcache.Item{Key: "a", Value: []byte("2")})
//
     it, err := mc.Get("foo")
	 if err != nil {
        panic(err)
    }
	fmt.Println(it)
// 	prototest()

	parse("localhost:33013", options, false)

}

// Replace this with your actual Options struct definition
type Options struct {
	idfa    string
	gaid    string
	adid    string
	dvid    string
	pattern string
}