# discovery-node
Draft implementation of `discovery-node` of the decentralized discovery protocol over Pss Swarm

### Run

#### Under TMUX window
```
bash run-tmux-demo.sh
```
This will launch four nodes inside a tmux window

#### Node by node
- Node0
```
go run *.go --config config0.yaml start
```

- Node1
```
go run *.go --config config1.yaml start
```
