#!/usr/bin/env python3

import argparse
import copy
import sys


def process_genesis(genesis, parsed_args):    
    genesis = genesis.replace('"uspolitics"', '"politics"')
    genesis = genesis.replace('"US Politics"', '"Politics"')
    genesis = genesis.replace(
        '"Discuss US politics like you would at Thanksgiving dinner', 
        '"Discuss politics like you would at Thanksgiving dinner'
    )
    
    return genesis


def init_default_argument_parser(prog_desc):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        '--exported-genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'),
        required=True,
    )
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    genesis = args.exported_genesis.read()

    print(process_genesis_func(genesis=genesis, parsed_args=args))

if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Migrate genesis.json uspolitics to politics',
    )
    main(parser, process_genesis)