package partitioner

import "sort"

// PartitionFn is a function that assigns a partition key to an environment variable.
type PartitionFn func(key, value string) string

// Partitioner splits an environment map into named partitions.
type Partitioner struct {
	fn PartitionFn
}

// New returns a new Partitioner using the provided partition function.
func New(fn PartitionFn) *Partitioner {
	return &Partitioner{fn: fn}
}

// Partition divides the given env map into named buckets based on the partition function.
func (p *Partitioner) Partition(env map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for k, v := range env {
		bucket := p.fn(k, v)
		if _, ok := result[bucket]; !ok {
			result[bucket] = make(map[string]string)
		}
		result[bucket][k] = v
	}
	return result
}

// PartitionNames returns a sorted list of partition names from the given result.
func PartitionNames(partitions map[string]map[string]string) []string {
	names := make([]string, 0, len(partitions))
	for name := range partitions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ByValueLength partitions variables into "short", "medium", or "long" based on value length.
func ByValueLength(short, long int) PartitionFn {
	return func(_, value string) string {
		switch {
		case len(value) <= short:
			return "short"
		case len(value) <= long:
			return "medium"
		default:
			return "long"
		}
	}
}
