package fleek

import (
  "context"
  "net/http"
  "time"

  "github.com/hasura/go-graphql-client"
  "golang.org/x/oauth2"
)

type FleekClient struct {
  client                      *graphql.Client
  httpClient                  *http.Client
  staticToken                 oauth2.TokenSource
  token                       string
}

type Team struct {
  Id                          string `json:"id"`
  Name                        string `json:"name"`
}

type EnvironmentVariable struct {
  Name                        string `json:"name"`
  Value                       string `json:"value"`
}

type BuildSettings struct {
  BuildCommand                string `json:"buildCommand"`
  BaseDirectoryPath           string `json:"baseDirectoryPath"`
  PublishDirectoryPath        string `json:"publishDirectoryPath"`
  DockerImage                 string `json:"dockerImage"`
  EnvironmentVariables      []EnvironmentVariable `json:"environmentVariables"`
}

type Source struct {
  // IPFS Source
  CID                         string `json:"cid,omitempty"`
  // Repository
  Type                        string `json:"type,omitempty"`
  URL                         string `json:"url,omitempty"`
  Branch                      string `json:"branch,omitempty"`
}

type DeploySettings struct {
  AutoPublishing              bool `json:"autoPublishing"`
  PRDeployPreviews            bool `json:"prDeployPreviews"`
  DfinityUseProxy             bool `json:"dfinityUseProxy"`
  Source                      Source `json:"source"`
}

type Repository struct {
  Commit                      string `json:"commit"`
  Branch                      string `json:"branch"`
  Owner                       string `json:"owner"`
  Name                        string `json:"name"`
  Message                     string `json:"message"`
}

type PublishedDeploy struct {
  Id                          interface{} `json:"id"`
  Status                      string `json:"status"`
  IPFSHash                    string `json:"ipfsHash"`
  Log                         string `json:"log"`
  PreviewImage                string `json:"previewImage"`
  AutoPublish                 bool   `json:"autoPublish"`
  Published                   bool   `json:"published"`
  Repository                  Repository `json:"repository"`
  TotalTime                   int    `json:"totalTime"`
  StartedAt                   time.Time `json:"startedAt"`
  CompletedAt                 time.Time `json:"completedAt"`
}

type Site struct {
  Id                          interface{} `json:"id"`
  Name                        string `json:"name"`
  Slug                        string `json:"slug"`
  Description                 string `json:"description"`
  Platform                    string `json:"platform"`

  Team                        Team `json:"team"`

  BuildSettings               BuildSettings   `json:"buildSettings"`
  DeploySettings              DeploySettings  `json:"deploySettings"`
  PublishedDeploy             PublishedDeploy `json:"publishedDeploy"`

  CreatedBy                   interface{} `json:"createdBy"`
  CreatedAt                   time.Time `json:"createdAt"`
  UpdatedAt                   time.Time `json:"updatedAt"`
}

type GraphqlSite struct {
  Id                          graphql.ID
  Name                        graphql.String
  Slug                        graphql.String
  Description                 graphql.String
  Platform                    graphql.String

  Team struct {
    Id                        graphql.ID
    Name                      graphql.String
  }

  BuildSettings struct {
    BuildCommand              graphql.String
    BaseDirectoryPath         graphql.String
    PublishDirectoryPath      graphql.String
    DockerImage               graphql.String
    EnvironmentVariables []struct {
      Name                    graphql.String
      Value                   graphql.String
    }
  }

  DeploySettings struct {
    AutoPublishing            graphql.Boolean
    PRDeployPreviews          graphql.Boolean
    DfinityUseProxy           graphql.Boolean
    Source struct {
      IPFSSource struct {
        CID                   graphql.String
      } `graphql:"... on IpfsSource"`
      Repository struct {
        Type                  graphql.String
        URL                   graphql.String
        Branch                graphql.String
      } `graphql:"... on Repository"`
    }
  }

  PublishedDeploy struct {
    Id                        graphql.ID
    Status                    graphql.String
    IpfsHash                  graphql.String
    PreviewImage              graphql.String
    AutoPublish               graphql.Boolean
    Published                 graphql.Boolean
    Log                       graphql.String
    Repository struct {
      Commit                  graphql.String
      Branch                  graphql.String
      Owner                   graphql.String
      Name                    graphql.String
      Message                 graphql.String
    }
    TotalTime                 graphql.Int
    StartedAt                 graphql.String
    CompletedAt               graphql.String
  }

  CreatedBy                   graphql.ID
  CreatedAt                   graphql.String
  UpdatedAt                   graphql.String
}

func (f *FleekClient) convertGraphqlSiteToSite(gqlSite GraphqlSite) (Site, error) {
  // BuildSettings
  buildSettings := BuildSettings{
    BuildCommand: string(gqlSite.BuildSettings.BuildCommand),
    BaseDirectoryPath: string(gqlSite.BuildSettings.BaseDirectoryPath),
    PublishDirectoryPath: string(gqlSite.BuildSettings.BaseDirectoryPath),
    DockerImage: string(gqlSite.BuildSettings.DockerImage),
  }
  for _, gqlEnvVar := range gqlSite.BuildSettings.EnvironmentVariables {
      buildSettings.EnvironmentVariables = append(
        buildSettings.EnvironmentVariables,
        EnvironmentVariable{
          Name: string(gqlEnvVar.Name),
          Value: string(gqlEnvVar.Value),
        },
      )
  }

  // DeploySettings
  deploySettings := DeploySettings{
    AutoPublishing: bool(gqlSite.DeploySettings.AutoPublishing),
    PRDeployPreviews: bool(gqlSite.DeploySettings.PRDeployPreviews),
    DfinityUseProxy: bool(gqlSite.DeploySettings.DfinityUseProxy),
    Source: Source{},
  }
  if gqlSite.DeploySettings.Source.IPFSSource.CID != "" {
    deploySettings.Source.CID =
      string(gqlSite.DeploySettings.Source.IPFSSource.CID)
  } else if gqlSite.DeploySettings.Source.Repository.URL != "" {
    deploySettings.Source.Type =
      string(gqlSite.DeploySettings.Source.Repository.Type)
    deploySettings.Source.URL =
      string(gqlSite.DeploySettings.Source.Repository.URL)
    deploySettings.Source.Branch =
      string(gqlSite.DeploySettings.Source.Repository.Branch)
  }

  // Site
  startedAt, _ := time.Parse(time.RFC3339, string(gqlSite.PublishedDeploy.StartedAt))
  completedAt, _ := time.Parse(time.RFC3339, string(gqlSite.PublishedDeploy.CompletedAt))

  createdAt, _ := time.Parse(time.RFC3339, string(gqlSite.CreatedAt))
  updatedAt, _ := time.Parse(time.RFC3339, string(gqlSite.UpdatedAt))

  site := Site{
    Id:       gqlSite.Id,
    Name:     string(gqlSite.Name),
    Slug:     string(gqlSite.Slug),
    Description: string(gqlSite.Description),
    Platform: string(gqlSite.Platform),

    Team: Team{
      Id:   gqlSite.Team.Id.(string),
      Name: string(gqlSite.Name),
    },

    BuildSettings: buildSettings,

    DeploySettings: deploySettings,

    PublishedDeploy: PublishedDeploy{
      Id:          gqlSite.PublishedDeploy.Id,
      Status:      string(gqlSite.PublishedDeploy.Status),
      IPFSHash:    string(gqlSite.PublishedDeploy.IpfsHash),
      PreviewImage: string(gqlSite.PublishedDeploy.PreviewImage),
      AutoPublish: bool(gqlSite.PublishedDeploy.AutoPublish),
      Published:   bool(gqlSite.PublishedDeploy.Published),
      Log:         string(gqlSite.PublishedDeploy.Log),
      Repository:  Repository{
        Commit: string(gqlSite.PublishedDeploy.Repository.Commit),
        Branch: string(gqlSite.PublishedDeploy.Repository.Branch),
        Owner: string(gqlSite.PublishedDeploy.Repository.Owner),
        Name: string(gqlSite.PublishedDeploy.Repository.Name),
        Message: string(gqlSite.PublishedDeploy.Repository.Message),
      },
      TotalTime: int(gqlSite.PublishedDeploy.TotalTime),
      StartedAt: startedAt,
      CompletedAt: completedAt,
    },

    CreatedBy: gqlSite.CreatedBy,
    CreatedAt: createdAt,
    UpdatedAt: updatedAt,
  }

  return site, nil
}


// Initializes new FleekClient
func New(token string) (*FleekClient, error) {
  fleekClient := new(FleekClient)

  fleekClient.token = token
  fleekClient.staticToken = oauth2.StaticTokenSource(&oauth2.Token{
    AccessToken: fleekClient.token,
    TokenType: " ",
  })
  fleekClient.httpClient = oauth2.NewClient(
    context.Background(),
    fleekClient.staticToken,
  )
  fleekClient.client = graphql.NewClient(
    "https://api.fleek.co/graphql",
    fleekClient.httpClient,
  )

  return fleekClient, nil
}

// Gets all sites of a team, with the team ID being the *slug* (e.g. `my-team`)
func (f *FleekClient) GetSitesByTeamId(teamId string) ([]Site, error) {
  var query struct {
      GetSitesByTeam struct {
        Sites                   []GraphqlSite
        NextToken               graphql.String
      } `graphql:"getSitesByTeam(teamId: $teamId, limit: 100)"`
  }

  vars := map[string]interface{} {
    "teamId": graphql.ID(teamId),
  }

  err := f.client.Query(context.Background(), &query, vars)
  if err != nil {
    return []Site{}, err
  }

  var sites []Site
  for _, querySite := range query.GetSitesByTeam.Sites {
    site, err := f.convertGraphqlSiteToSite(querySite)
    if err != nil {
      return sites, err
    }
    sites = append(sites, site)
  }

  return sites, nil
}

// Gets a single site by its slug
func (f *FleekClient) GetSiteBySlug(slug string) (Site, error) {
  var query struct {
      GetSiteBySlug struct {
       GraphqlSite
      } `graphql:"getSiteBySlug(slug: $slug)"`
  }

  vars := map[string]interface{} {
    "slug": graphql.String(slug),
  }

  err := f.client.Query(context.Background(), &query, vars)
  if err != nil {
    return Site{}, err
  }

  site, err := f.convertGraphqlSiteToSite(query.GetSiteBySlug.GraphqlSite)
  return site, err
}

