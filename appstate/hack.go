package appstate

//// PatchInfo contains information about a patch to the app state.
//// A patch can contain multiple mutations, as long as all mutations are in the same app state type.
//type PatchInfo struct {
//	// Timestamp is the time when the patch was created. This will be filled automatically in EncodePatch if it's zero.
//	Timestamp time.Time
//	// Type is the app state type being mutated.
//	Type WAPatchName
//	// Mutations contains the individual mutations to apply to the app state in this patch.
//	Mutations []MutationInfo
//	// Operation is SET / REMOVE
//	Operation waServerSync.SyncdMutation_SyncdOperation
//}
//
//func BuildContact(target types.JID, fullName string) PatchInfo {
//	return PatchInfo{
//		Type:      WAPatchCriticalUnblockLow,
//		Operation: waServerSync.SyncdMutation_SET,
//		Mutations: []MutationInfo{{
//			Index:   []string{IndexContact, target.String()},
//			Version: 2,
//			Value: &waSyncAction.SyncActionValue{
//				ContactAction: &waSyncAction.ContactAction{
//					FullName:                 &fullName,
//					SaveOnPrimaryAddressbook: proto.Bool(true),
//				},
//			},
//		}},
//	}
//}
//
//func RemoveContact(target types.JID) PatchInfo {
//	return PatchInfo{
//		Type:      WAPatchCriticalUnblockLow,
//		Operation: waServerSync.SyncdMutation_REMOVE,
//		Mutations: []MutationInfo{{
//			Index:   []string{IndexContact, target.String()},
//			Version: 2,
//			Value: &waSyncAction.SyncActionValue{
//				ContactAction: &waSyncAction.ContactAction{},
//			},
//		}},
//	}
//}
