package main

import (
	"fmt"
	"strconv"
	"strings"
	"google.golang.org/protobuf/proto"
	proto_struct "github.com/balabanas/memc-protobuf-go/proto"
)

func prototest() {
    // Test if we can successfully marshall and unmarshall our messages with protobuf library,
    // without actually doing some export to memcached
	sample := "idfa\t1rfw452y52g2gq4g\t55.55\t42.42\t1423,43,567,3,7,23\ngaid\t7rfw452y52g2gq4g\t55.55\t42.42\t7423,424"
	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
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
			fmt.Println("Packed/unpacked protobuf objects are equal.")
		} else {
			fmt.Println("Packed/unpacked protobuf objects are NOT equal.")
		}
	}
}