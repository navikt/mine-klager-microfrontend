export const isDeployedToProd = process.env.NAIS_CLUSTER_NAME === 'prod-gcp';
export const isDeployedToDev = process.env.NAIS_CLUSTER_NAME === 'dev-gcp';
export const isDeployed = isDeployedToProd || isDeployedToDev;
export const isLocal = !isDeployed;
