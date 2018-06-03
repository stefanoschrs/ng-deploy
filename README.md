# Angular & AngularJS deployer

## Usage
```
cd my-app
ng-deploy --env dev --build --sync
```
- **env** string  
Select environment
- **build**  
Build project
- **sync**  
Sync files
- **sync-type**  
How to sync files: ssh, aws
- **backup**  
Backup target

## Configuration
- Create a `.ng-deploy.json` file in your project dir with your **env** configurations (Example [here](.ng-deploy.example.json))
- Change your `package.json` accordingly to contain the **build:\<env name\>** scripts:
```
{
  "name": "my-app",
  "scripts": {
    ...
    "build:prod": "ng build --prod --aot --build-optimizer --base-href='https://www.google.com'",
    "build:dev": "ng build --base-href='https://dev.google.com'",
    ...
  }
}

```