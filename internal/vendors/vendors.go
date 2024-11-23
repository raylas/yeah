package vendors

import "strings"

type VendorEntry struct {
	Oui     string
	Name    string
	Address string
}

type VendorNode struct {
	Children map[rune]*VendorNode
	IsEnd    bool
	Data     *VendorEntry
}

type Vendors struct {
	Root *VendorNode
	// Revision string
}

func New() *Vendors {
	return &Vendors{
		Root: &VendorNode{
			Children: make(map[rune]*VendorNode),
		},
	}
}

func (v *Vendors) Insert(prefix string, data *VendorEntry) {
	node := v.Root
	for _, char := range prefix {
		if _, found := node.Children[char]; !found {
			node.Children[char] = &VendorNode{Children: make(map[rune]*VendorNode)}
		}
		node = node.Children[char]
	}
	node.IsEnd = true
	node.Data = data
}

func (v *Vendors) Search(prefix string) []*VendorEntry {
	normalized := normalize(prefix)

	node := v.Root
	var lastMatchNode *VendorNode

	for i := 0; i < len(normalized); i++ {
		char := rune(normalized[i])
		if next, found := node.Children[char]; found {
			node = next
			if node.IsEnd {
				lastMatchNode = node
			}
		} else {
			// // If we can't match the full string but we're at a valid node,
			// // return all matches from current position
			// if i < 8 { // Less than 8 chars means it's a partial search
			// 	return v.collect(node)
			// }
			break
		}
	}

	// For full MAC addresses, return the longest match we found
	if len(normalized) >= 8 && lastMatchNode != nil {
		return []*VendorEntry{lastMatchNode.Data}
	}

	// For partial searches, return all matches from where we ended
	return v.collect(node)
}

func (v *Vendors) collect(node *VendorNode) []*VendorEntry {
	var entries []*VendorEntry
	if node.IsEnd {
		entries = append(entries, node.Data)
	}
	for _, child := range node.Children {
		entries = append(entries, v.collect(child)...)
	}
	return entries
}

func normalize(prefix string) string {
	replacer := strings.NewReplacer(":", "", "-", "", ".", "")
	replaced := replacer.Replace(prefix)
	return strings.ToUpper(replaced)
}
