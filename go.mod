module github.com/audactai/audact-go-sdk

go 1.24

retract v0.3.0 // Calls.List, Calls.Get, Evidence.Reconstruct and Policy.Get fail to decode real API responses; use v0.4.0.
