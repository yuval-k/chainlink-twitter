#!/usr/bin/env node
var MarkdownIt = require('markdown-it'),
md = new MarkdownIt();
const fs = require('fs'); 
const { type } = require('os');

var text = fs.readFileSync('./setup_local_testnet.md', {encoding:'utf8', flag:'r'});

let output = console.log;

// output prefix:
output("#!/bin/bash -e\n");

var tree = md.parse(text, {});
let heading = "";
let new_section = false;
for (let i = 0 ; i < tree.length; i++) {
    let token = tree[i];
    // current section
    if (token.type == "heading_open" && token.tag == "h1"){
        // get content of next token
        let nextoken = tree[i+1];
        heading = tree[i+1].content;
        new_section = true;
    }
    if (token.tag == "code" && token.info == 'bash'){
        // if we ar enot in an optional section
        let optional = heading.toLowerCase().includes("optional");
        let prereq = heading.toLowerCase().includes("prerequisites");
        if (!optional && ! prereq) {
            if (new_section){
                output("# section: " + heading);
                new_section = false;
            }
            output(token.content);
        }
    }
}

// output suffix:
output("echo export ORACLE_ADDR=$ORACLE_ADDR\necho export TWITTER_JOB_ID=$TWITTER_JOB_ID\n");