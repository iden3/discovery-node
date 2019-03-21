#!/bin/sh

SESSION='discovery-node'

tmux new-session -d -s $SESSION
tmux split-window -d -t 0 -v
tmux split-window -d -t 0 -h
tmux split-window -d -t 2 -h

tmux send-keys -t 0 'go run *.go --config config0.yaml start' enter
tmux send-keys -t 1 'go run *.go --config config1.yaml start' enter
tmux send-keys -t 2 'go run *.go --config config2.yaml start' enter
tmux send-keys -t 3 'go run *.go --config config3.yaml start' enter

tmux attach
