package main

import (
"testing"
)

func BenchmarkOriginalMode(b *testing.B) {
lines := []string{"0,0", "10,0", "10,10", "0,10"}
b.ResetTimer()
for i := 0; i < b.N; i++ {
processCoordinates(lines, "original")
}
}

func BenchmarkContainedModeSmall(b *testing.B) {
lines := []string{"0,0", "10,0", "10,10", "0,10"}
b.ResetTimer()
for i := 0; i < b.N; i++ {
processCoordinates(lines, "contained")
}
}

func BenchmarkContainedModeMedium(b *testing.B) {
lines := []string{"0,0", "1000,0", "1000,1000", "0,1000"}
b.ResetTimer()
for i := 0; i < b.N; i++ {
processCoordinates(lines, "contained")
}
}

func BenchmarkContainedModeLarge(b *testing.B) {
lines := []string{"0,0", "10000,0", "10000,10000", "0,10000"}
b.ResetTimer()
for i := 0; i < b.N; i++ {
processCoordinates(lines, "contained")
}
}

func BenchmarkContainedModeVeryLarge(b *testing.B) {
lines := []string{"0,0", "100000,0", "100000,100000", "0,100000"}
b.ResetTimer()
for i := 0; i < b.N; i++ {
processCoordinates(lines, "contained")
}
}
