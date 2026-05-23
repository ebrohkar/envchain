// Package grouper partitions environment variable maps into named groups
// using a configurable classification function.
//
// A GroupFn maps each key to a group name. Keys that return an empty
// string are placed in the ungrouped bucket.
//
// Example usage:
//
//	g := grouper.New(grouper.ByPrefix("_"))
//	groups := g.Group(env)
//	for _, name := range grouper.GroupNames(groups) {
//		fmt.Println(name, groups[name])
//	}
package grouper
