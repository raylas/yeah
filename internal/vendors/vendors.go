package vendors

import "strings"

type VendorEntry struct {
	Oui     string `json:"oui"`
	Name    string `json:"name"`
	Address string `json:"address"`
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
	normalized := normalize(prefix)
	node := v.Root

	for _, char := range normalized {
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
	var lastMatch *VendorNode

	for _, char := range normalized {
		next, found := node.Children[char]
		if !found {
			// For full MAC addresses, return the last valid OUI match if we found one
			if len(normalized) >= 12 && lastMatch != nil {
				return []*VendorEntry{lastMatch.Data}
			}
			return nil
		}
		node = next
		if node.IsEnd {
			lastMatch = node
		}
	}

	// For full MAC addresses, return only the longest match
	if len(normalized) >= 12 && lastMatch != nil {
		return []*VendorEntry{lastMatch.Data}
	}

	// For partial searches, return all matches from where we ended up
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
	if len(replaced) > 12 {
		replaced = replaced[:12]
	}
	return strings.ToUpper(replaced)
}
