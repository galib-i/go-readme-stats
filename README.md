# Go README Stats
Add dynamic language usage stats from your projects directly to your personal README as a minimal GitHub-style card.

> [!IMPORTANT]
> This will only work for your own GitHub account. It requires tokens from both **GitHub** and **Vercel**.
> 
> *A Vercel account is required for this setup: the free tier is sufficient.*

> [!CAUTION]
> Do not share your GitHub or Vercel tokens with anyone, as they provide access to your account.

## Setup
### Steps
1. Clone or [fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) this repository
2. Sign in to [Vercel](https://vercel.com), go to the dashboard and click **Add New > Project**;
under **Import Git Repository**, select `go-readme-stats`
3. Create a [GitHub Personal Access Token (PAT)](https://github.com/settings/tokens/new) with **repo** and **read:user** scopes  
   *Set the token expiry as you see appropriate*
4. In the **Environment Variables** section on Vercel, set:  
   - **Name:** `GITHUB_TOKEN`  
   - **Value:** *your GitHub token*
  
   Then press **Deploy**

If you don’t want to use GitHub Actions, you’ve finished the setup. 

To use GitHub Actions for a simple workflow that updates language colours, dependencies and tests before deployment:

5. In your `go-readme-stats` repository, go to **Settings > Secrets and variables > Actions** and press **New repository secret**:  
   - Copy the **PROJECT_ID** from your Vercel project settings and create a secret in GitHub:  
     - **Name:** `VERCEL_PROJECT_ID`  
     - **Value:** *your Vercel `PROJECT_ID`*  
   - Copy the **TEAM_ID** from your Vercel dashboard ([find your team ID](https://vercel.com/docs/accounts#find-your-team-id)) and create another GitHub secret:  
     - **Name:** `VERCEL_ORG_ID`  
     - **Value:** *your Vercel `TEAM_ID`*  
   - Generate a [Vercel account token](https://vercel.com/guides/how-do-i-use-a-vercel-api-access-token) and create a final GitHub secret:  
     - **Name:** `VERCEL_TOKEN`  
     - **Value:** *your Vercel token*

    *Alternatively, find `projectId` and `orgId` in `.vercel/project.json`*

6. Enable all GitHub Action workflows on your repository  

## Usage
The base endpoint is `/api/langs`.

By default, up to 6 languages are shown. If there are more, the top 5 are shown individually and the rest are grouped under *Other*.

### Basic Example
Add this to your `README.md`:
```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs)
```
### Themes
- Dark (default)
```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?theme=dark)
```
<img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/10df753f-cb46-4448-b5bb-ebcb4f1b62a6" />  
<br><br>

- Soft-dark
```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?theme=soft-dark)
```
<img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/a2613437-4a57-4c35-ab39-a9d01214faa1" />  
<br><br>

- Light
```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?theme=light)
```
<img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/3ef24020-88b0-420b-a256-4bb2ad8c1b5e" />  
<br><br>

>[!TIP]
>Wrap the card with GitHub’s fragment identifiers to show cards conditionally in light or dark mode (`#gh-light-mode-only` / `#gh-dark-mode-only`):
```yaml
[![Light Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?theme=light)](https://github.com/<YOUR USERNAME>/go-readme-stats#gh-light-mode-only)
[![Dark Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs)](https://github.com/<YOUR USERNAME>/go-readme-stats#gh-dark-mode-only)
```
### Custom Header
Change the title in the card with the `?header=` parameter (use %20 for spaces):

```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?header=Used%20Languages)
```

### Percentage Calculation Method
By default, percentages are calculated as *language bytes ÷ total bytes*.

This can overweight languages that use more bytes to express the same logic, so you can switch to geometric mean with `?mode=geometric`, which downweights extremes and attempts to better reflect the typical language mix across projects:

```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?mode=geometric)
```
<img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/10df753f-cb46-4448-b5bb-ebcb4f1b62a6" /> <img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/13c5de13-68de-4423-9e7e-e2cc800f4340" />  


### Combine Options
Combine parameters with `&`:
```yaml
![Languages Card](https://<YOUR VERCEL DOMAIN>/api/langs?theme=light&mode=geometric)
```
<br><br>
*Example in my [profile](https://github.com/galib-i)*
