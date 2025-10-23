#!/bin/sh
go clean -testcache 
modulestest=`find . -mindepth 1 -type d `

# color definition
red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
blue=`tput setaf 4`
magenta=`tput setaf 5`
cyan=`tput setaf 6`
reset=`tput sgr0`
bold=`tput bold`

# TEST EACH MODULE, fail as soon as one fails
for module in  ${modulestest}; do
    # test if it is a valid code container 
    testmodule=$module"_test"
    # test if _test file exists: if not, not a golang code directory
    if [ ! -d $testmodule ] ; then
        continue 
    fi
    # test if it contains go files
    counter=`cat $module/*.go | wc -l`
    if [ $counter = 0 ]; then 
        continue
    fi 

    # assuming it is a go source directory
    # Run tests, and get result 
    testresult=`go test ./$testmodule/ | awk '{print $1 "  " $3}'`
    teststatus=`echo $testresult | awk '{print $1}'`
    testtime=`echo $testresult | awk '{print $2}'`
    # if success, print status and time 
    # else rerun test, no output capture
    if [ $teststatus = "ok" ]; then
        echo -n "${bold}$module  ${reset}"
        echo -n "${green}${bold}[GO]${reset} in $testtime "
        echo -n "${yellow}    code:${reset} "
        counter=`cat $module/*.go | wc -l`
        echo -n "${yellow}$counter${reset} "
        testcounter=`cat $testmodule/*.go | wc -l`
        echo -n "${yellow}    test: ${reset}"
        echo -n "${yellow}$testcounter ${reset}"
        echo
    else
        echo -n "${bold}$module  ${reset}"
        echo -n "${red}${bold}[FAILING]  ${reset}"
        echo
        go test ./$testmodule/
        exit 1
    fi
done 
