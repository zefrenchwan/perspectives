#!/bin/sh
go clean -testcache 
modulestest="commons"

# color definition
red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
blue=`tput setaf 4`
magenta=`tput setaf 5`
cyan=`tput setaf 6`
reset=`tput sgr0`

# TEST EACH MODULE, fail as soon as one fails
for module in  ${modulestest}; do 
    # print module info
    echo -n "${reset}$module ${reset}"
    echo -n "${yellow}    code:${reset} "
    counter=`cat $module/*.go | wc -l`
    echo -n "${yellow}$counter${reset} "
    testmodule=$module"_test"
    testcounter=`cat $testmodule/*.go | wc -l`
    echo -n "${yellow}    test: ${reset}"
    echo -n "${yellow}$testcounter ${reset}"
    # Run tests, and get result 
    testresult=`go test ./$testmodule/ | awk '{print $1 "  " $3}'`
    teststatus=`echo $testresult | awk '{print $1}'`
    testtime=`echo $testresult | awk '{print $2}'`
    # if success, print status and time 
    # else rerun test, no output capture
    if [ $teststatus = "ok" ]; then
        echo -n "${green}    status: OK ( $testtime )${reset}"
        echo
    else
        echo -n "${red}    status: FAILURE${reset}"
        echo
        go test ./$testmodule/
        exit -1
    fi
done 
