package abe

import (
)

type ABitOfEverythingNested struct {
    Name  string  `json:"name,omitempty"`
    Amount  int64  `json:"amount,omitempty"`
    Ok  NestedDeepEnum  `json:"ok,omitempty"`
    
}
