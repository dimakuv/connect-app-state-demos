#!/usr/bin/env python3

# Copyright (C) 2024 Intel Labs
# SPDX-License-Identifier: BSD-3-Clause

import argparse
import os
import pickle
import signal
import sys
import time

RESTORE_FILENAME_ENVVAR = "RESTORE_FROM_FILE"

person = { } # set below

def checkpoint_before_exit(person, filename):
    with open(filename, 'wb') as file:
        pickle.dump(person, file)
        print('Checkpointed into', filename, flush=True)

def restore_before_start(filename):
    with open(filename, 'rb') as file:
        person = pickle.load(file)
        print('Restored from', filename, flush=True)
        return person

def handler(signum, frame):
    global args
    global person
    print('Signal handler called with signal', signum, flush=True)
    checkpoint_before_exit(person, args.file)
    sys.exit(0)

def main_loop(person):
    while True:
        print(person, flush=True)
        person['age'] += 1
        time.sleep(1)

parser = argparse.ArgumentParser()
parser.add_argument('--name', help='Person\'s name', default='John Doe', type=str)
parser.add_argument('-a', '--age', help='Person\'s initial age', default=0, type=int)
parser.add_argument('--file', help='File to dump the state to', default='dump.dat', type=str)
args = parser.parse_args()

if RESTORE_FILENAME_ENVVAR in os.environ:
    person = restore_before_start(os.environ[RESTORE_FILENAME_ENVVAR])
else:
    person = { 'name': args.name, 'age': args.age }

signal.signal(signal.SIGINT, handler)
signal.signal(signal.SIGTERM, handler)
main_loop(person)
