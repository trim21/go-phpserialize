module bench

go 1.19

require (
	github.com/elliotchance/phpserialize v1.4.0
	github.com/trim21/go-phpserialize v0.0.0
)

require go4.org/unsafe/assume-no-moving-gc v0.0.0-20231121144256-b99613f794b6 // indirect

replace github.com/trim21/go-phpserialize => ../
