import Amplify from "aws-amplify";

const AWS_CONFIGURATION=''

const AmplifyService={
    Amplify.configure(
        {
            Auth:{
                identityPoolId: process.env.VUE_APP_AMPLIFY_AUTH_IDENTITY_POOL_ID,
                userPoolId: process.env.VUE_APP_AMPLIFY_AUTH_USER_POOL_ID,
                userPoolWebClientId: process.env.VUE_APP_AMPLIFY_AUTH_USER_POOL_WEB_CLIENT_ID,
                identityPoolRegion: process.env.VUE_APP_DEPLOY_AWS_REGION,
                region: process.env.VUE_APP_AMPLIFY_AUTH_REGION
            },
            API:{
                endpoints: [
                    {
                      name: PlatformAdminApi,
                      endpoint: endpointUrl,
                      region: process.env.VUE_APP_AMPLIFY_AUTH_REGION,
                      custom_header: async () => {
                          try {
                              return {
                                  Authorization: `Bearer ${(await Auth.currentSession()).getIdToken().getJwtToken()}`
                              };
                          } catch (e) {
                              return {};
                          }
                      }
                    }
                  ]
            }
        }
    )

    // import Amplify from "@aws-amplify/core";
    // import Auth from "@aws-amplify/auth";
    // import 
}

export {
    AmplifyService
}