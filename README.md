# Kusto Dashboard Sync Tool

## **Warning:** This tool is untested and very hackly, back up all dashboards before trying out it out.

# Install
`kusto-dashboards-sync` is a command line tool to fetch or publish a dataexplorer dashboard.

To install the tool run:
```
 go install github.com/omeshp/kusto-dashboards-sync@latest
```

# Setup
- Access token for authentication to data explorer is stored locally in `.env` file and can be setup using below az commands:
```
az login
echo ACCESS_TOKEN=$(az account get-access-token --resource "35e917a9-4d95-4062-9d97-5781291353b9" --query "accessToken") > .env
```

- Setup `config.yml` having id for the dashboard to sync:
```
echo "dashboard_id: af4de11d-baac-4bd4-b2fd-e979f68f31be" > config.yml
```

- Add `.env` to `.gitignore` if you plan on syncing dashboards to github.
```
echo "\n.env" >> .gitignore      
```

# Usage
`kusto-dashboards-sync` provides below command line options:

- Export dashboard which places extracts queries to `queries` folder, with includes in `dashboard.yml`

```
kusto-dashboards-sync pull [dashboard id]
```

- Will process template `dashboard.yml` and push updates to the dashboard

```
kusto-dashboards-sync push
```

If no dashboard id is specified the dashboard to pull/push is picked from `config.yml` file.
Example `config.yml`:
```
dashboard_id: af4de11d-baac-4bd4-b2fd-e979f68f31be
```

