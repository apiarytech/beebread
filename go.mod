module github.com/apiarytech/beebread

go 1.27.1

require github.com/apiarytech/royaljelly v0.0.5-beta1

// The iec package beebread is written against is not in a tagged royaljelly
// release yet; remove this once one is published.
replace github.com/apiarytech/royaljelly => ../royaljelly
