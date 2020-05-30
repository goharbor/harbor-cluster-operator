package database

// Deploy reconcile will deploy database cluster if that does not exist.
// It does:
// - check postgre does exist
// - create any new postgresqls.acid.zalan.do CRs
// - create postgre connection secret
// It does not:
// - perform any postgresqls downscale (left for downscale phase)
// - perform any postgresqls upscale (left for upscale phase)
// - perform any pod upgrade (left for rolling upgrade phase)
func (postgre *PostgreSQLReconciler) Deploy() error {

	if postgre.HarborCluster.Spec.Database.Kind == "external" {
		return nil
	}

	return nil
}
