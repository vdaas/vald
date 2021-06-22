package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func CheckPayloadObjectAccess(m dsl.Matcher) {
	m.Import("github.com/vdaas/vald/apis/grpc/v1/payload")

	m.Match(
		`$*_ = $*_, $x.$y, $*_`,
		`$*_, $x.$y, $*_`,
		`$x.$y != $_`,
		`$x.$y == $_`,
	).Where(!m["y"].Text.Matches(`Get.+`) &&
		(m["x"].Type.Is(`*payload.Search_Request`) ||
			m["x"].Type.Is(`*payload.Search_MultiRequest`) ||
			m["x"].Type.Is(`*payload.Search_IDRequest`) ||
			m["x"].Type.Is(`*payload.Search_MultiIDRequest`) ||
			m["x"].Type.Is(`*payload.Search_ObjectRequest`) ||
			m["x"].Type.Is(`*payload.Search_MultiObjectRequest`) ||
			m["x"].Type.Is(`*payload.Search_Config`) ||
			m["x"].Type.Is(`*payload.Search_Response`) ||
			m["x"].Type.Is(`*payload.Search_Responses`) ||
			m["x"].Type.Is(`*payload.Search_StreamResponse`) ||
			m["x"].Type.Is(`*payload.Filter_Target`) ||
			m["x"].Type.Is(`*payload.Filter_Config`) ||
			m["x"].Type.Is(`*payload.Insert_Request`) ||
			m["x"].Type.Is(`*payload.Insert_MultiRequest`) ||
			m["x"].Type.Is(`*payload.Insert_ObjectRequest`) ||
			m["x"].Type.Is(`*payload.Insert_MultiObjectRequest`) ||
			m["x"].Type.Is(`*payload.Insert_Config`) ||
			m["x"].Type.Is(`*payload.Update_Request`) ||
			m["x"].Type.Is(`*payload.Update_MultiRequest`) ||
			m["x"].Type.Is(`*payload.Update_ObjectRequest`) ||
			m["x"].Type.Is(`*payload.Update_MultiObjectRequest`) ||
			m["x"].Type.Is(`*payload.Update_Config`) ||
			m["x"].Type.Is(`*payload.Upsert_Request`) ||
			m["x"].Type.Is(`*payload.Upsert_MultiRequest`) ||
			m["x"].Type.Is(`*payload.Upsert_ObjectRequest`) ||
			m["x"].Type.Is(`*payload.Upsert_MultiObjectRequest`) ||
			m["x"].Type.Is(`*payload.Upsert_Config`) ||
			m["x"].Type.Is(`*payload.Remove_Request`) ||
			m["x"].Type.Is(`*payload.Remove_MultiRequest`) ||
			m["x"].Type.Is(`*payload.Remove_Config`) ||
			m["x"].Type.Is(`*payload.Object_VectorRequest`) ||
			m["x"].Type.Is(`*payload.Object_Distance`) ||
			m["x"].Type.Is(`*payload.Object_StreamDistance`) ||
			m["x"].Type.Is(`*payload.Object_ID`) ||
			m["x"].Type.Is(`*payload.Object_IDs`) ||
			m["x"].Type.Is(`*payload.Object_Vector`) ||
			m["x"].Type.Is(`*payload.Object_Vectors`) ||
			m["x"].Type.Is(`*payload.Object_StreamVector`) ||
			m["x"].Type.Is(`*payload.Object_StreamVector_Vector`) ||
			m["x"].Type.Is(`*payload.Object_StreamVector_Status`) ||
			m["x"].Type.Is(`*payload.Object_Blob`) ||
			m["x"].Type.Is(`*payload.Object_StreamBlob`) ||
			m["x"].Type.Is(`*payload.Object_StreamBlob_Blob`) ||
			m["x"].Type.Is(`*payload.Object_StreamBlob_Status`) ||
			m["x"].Type.Is(`*payload.Object_Location`) ||
			m["x"].Type.Is(`*payload.Object_StreamLocation`) ||
			m["x"].Type.Is(`*payload.Object_StreamLocation_Location`) ||
			m["x"].Type.Is(`*payload.Object_StreamLocation_Status`) ||
			m["x"].Type.Is(`*payload.Object_Locations`) ||
			m["x"].Type.Is(`*payload.Meta_Key`) ||
			m["x"].Type.Is(`*payload.Meta_Keys`) ||
			m["x"].Type.Is(`*payload.Meta_Val`) ||
			m["x"].Type.Is(`*payload.Meta_Vals`) ||
			m["x"].Type.Is(`*payload.Meta_KeyVal`) ||
			m["x"].Type.Is(`*payload.Meta_KeyVals`) ||
			m["x"].Type.Is(`*payload.Control_CreateIndexRequest`) ||
			m["x"].Type.Is(`*payload.Replication_Recovery`) ||
			m["x"].Type.Is(`*payload.Replication_Rebalance`) ||
			m["x"].Type.Is(`*payload.Replication_Agents`) ||
			m["x"].Type.Is(`*payload.Discoverer_Request`) ||
			m["x"].Type.Is(`*payload.Backup_GetVector_Request`) ||
			m["x"].Type.Is(`*payload.Backup_GetVector_Owner`) ||
			m["x"].Type.Is(`*payload.Backup_Locations_Request`) ||
			m["x"].Type.Is(`*payload.Backup_Remove_Request`) ||
			m["x"].Type.Is(`*payload.Backup_Remove_RequestMulti`) ||
			m["x"].Type.Is(`*payload.Backup_IP_Register_Request`) ||
			m["x"].Type.Is(`*payload.Backup_IP_Remove_Request`) ||
			m["x"].Type.Is(`*payload.Backup_Vector`) ||
			m["x"].Type.Is(`*payload.Backup_Vectors`) ||
			m["x"].Type.Is(`*payload.Backup_Compressed_Vector`) ||
			m["x"].Type.Is(`*payload.Backup_Compressed_Vectors`) ||
			m["x"].Type.Is(`*payload.Info_Index_Count`) ||
			m["x"].Type.Is(`*payload.Info_Index_UUID_Committed`) ||
			m["x"].Type.Is(`*payload.Info_Index_UUID_Uncommitted`) ||
			m["x"].Type.Is(`*payload.Info_Pod`) ||
			m["x"].Type.Is(`*payload.Info_Node`) ||
			m["x"].Type.Is(`*payload.Info_CPU`) ||
			m["x"].Type.Is(`*payload.Info_Memory`) ||
			m["x"].Type.Is(`*payload.Info_Pods`) ||
			m["x"].Type.Is(`*payload.Info_Nodes`) ||
			m["x"].Type.Is(`*payload.Info_IPs`))).
		At(m["y"]).
		Report("Avoid to access struct fields directly").
		Suggest(`Get$y()`)
}

// from Ruleguard by example
// https://go-ruleguard.github.io/by-example/matching-comments
func PrintFmt(m dsl.Matcher) {
	// Here we scan a $s string for %s, %d and %v.
	m.Match(`fmt.Println($s, $*_)`,
		`fmt.Print($s, $*_)`).
		Where(m["s"].Text.Matches(`%[sdv]`)).
		Report("found formatting directive in non-formatting call")
}
