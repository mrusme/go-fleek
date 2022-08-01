package fleek

import (
  "context"
  "net/http"

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

type PublishedDeploy struct {
  Id                          interface{} `json:"id"`
  Status                      string `json:"status"`
  IPFSHash                    string `json:"ipfsHash"`
  Log                         string `json:"log"`
  CompletedAt                 string `json:"completedAt"`
}

type Site struct {
  Id                          interface{} `json:"id"`
  Name                        string `json:"name"`
  Platform                    string `json:"platform"`
  PublishedDeploy             PublishedDeploy `json:"publishedDeploy"`
  Team                        Team `json:"team"`
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
    site := Site{
      Id:       querySite.Id,
      Name:     string(querySite.Name),
      Platform: string(querySite.Platform),
      PublishedDeploy: PublishedDeploy{
        Id:          querySite.PublishedDeploy.Id,
        Status:      string(querySite.PublishedDeploy.Status),
        IPFSHash:    string(querySite.PublishedDeploy.IpfsHash),
        Log:         string(querySite.PublishedDeploy.Log),
        CompletedAt: string(querySite.PublishedDeploy.CompletedAt),
      },
      Team: Team{
        Id:   querySite.Team.Id.(string),
        Name: string(querySite.Name),
      },
    }
    sites = append(sites, site)
  }

  return sites, nil
}


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

  site := Site{
    Id:       query.GetSiteBySlug.Id,
    Name:     string(query.GetSiteBySlug.Name),
    Platform: string(query.GetSiteBySlug.Platform),
    PublishedDeploy: PublishedDeploy{
      Id:          query.GetSiteBySlug.PublishedDeploy.Id,
      Status:      string(query.GetSiteBySlug.PublishedDeploy.Status),
      IPFSHash:    string(query.GetSiteBySlug.PublishedDeploy.IpfsHash),
      Log:         string(query.GetSiteBySlug.PublishedDeploy.Log),
      CompletedAt: string(query.GetSiteBySlug.PublishedDeploy.CompletedAt),
    },
    Team: Team{
      Id:   query.GetSiteBySlug.Team.Id.(string),
      Name: string(query.GetSiteBySlug.Team.Name),
    },
  }

  return site, nil
}

