// Copyright (C) 2024 Intel Labs
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
    "bytes"
    "encoding/gob"
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

const RESTORE_FILENAME_ENVVAR = "RESTORE_FROM_FILE"

type Person struct {
    Name string
    Age  int
}

func checkpoint_before_exit(personPtr *Person, filename string) {
    var b bytes.Buffer
    enc := gob.NewEncoder(&b)

    if err := enc.Encode(personPtr); err != nil {
        log.Fatal(err)
    }
    if err := os.WriteFile(filename, b.Bytes(), 0644); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Checkpointed into", filename)
}

func restore_before_start(personPtr *Person, filename string) {
    dat, err := os.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    b := bytes.NewBuffer(dat)
    dec := gob.NewDecoder(b)

    if err := dec.Decode(personPtr); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Restored from", filename)
}

func main_loop(wg *sync.WaitGroup, person *Person, stop_chan chan struct{}) {
    defer wg.Done()

    for {
        select {
        case <-stop_chan:
            return
        default:
        }

        fmt.Println(*person)
        person.Age += 1
        time.Sleep(time.Second)
    }
}

func main() {
    var person Person

    namePtr := flag.String("name", "John Doe", "Person's name")
    agePtr  := flag.Int("age", 0, "Person's initial age")
    dump_filenamePtr := flag.String("file", "dump.dat", "File to dump the state to")
    flag.Parse()

    restore_filename, do_restore := os.LookupEnv(RESTORE_FILENAME_ENVVAR)
    if do_restore {
        restore_before_start(&person, restore_filename)
    } else {
        person = Person{Name: *namePtr, Age: *agePtr}
    }

    stop_chan := make(chan struct{})

    signal_chan := make(chan os.Signal)
    signal.Notify(signal_chan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-signal_chan
        close(stop_chan)
    }()

    wg := sync.WaitGroup{}
    wg.Add(1)
    go main_loop(&wg, &person, stop_chan)
    wg.Wait()

    checkpoint_before_exit(&person, *dump_filenamePtr)
}
