#!/bin/sh

SESSION='discovery-node'

tmux new-session -d -s $SESSION
tmux split-window -d -t 0 -v
tmux split-window -d -t 0 -h
tmux split-window -d -t 2 -h

tmux send-keys -t 0 'go run main.go --config config0.yaml start' enter
tmux send-keys -t 1 'go run main.go --config config1.yaml start' enter
tmux send-keys -t 2 'go run main.go --config config2.yaml start' enter

tmux select-pane -t 3
tmux send-keys -t 3 'curl http://127.0.0.1:4000/id/0x47a2b2353f1a55e4c975b742a7323c027160b4e3 | jq' enter
tmux send-keys -t 3 '' enter
tmux send-keys -t 3 'python test_discovery.py'

tmux attach
