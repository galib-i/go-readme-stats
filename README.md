# Go README Stats
Add language usage stats from your projects directly to your personal README as a GitHub-style card. 

<p align="center">
  <img alt="Languages Card" src="https://raw.githubusercontent.com/galib-i/go-readme-stats/output/languages-dark.svg">
</p>

> [!IMPORTANT]
> This will only work for your own GitHub account: it uses a Personal Access Token (PAT).

- Use `stats.yml` or flags to tweak your card's theme, header, calculation mode, maximum shown languages and ignore list.
- An updated card is generated every 12 hours (and whenever `stats.yml` changes).
- Colour updates are fetched on the 1st of every month.
- Dependencies are updated on the 1st every third month.

## Get Started
Cards are generated using GitHub Actions and added to an orphan branch, keeping the main branch free of automated update commits.

1. Clone or [fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) this repository, enabling all GitHub Action workflows.
2. Create a [GitHub Personal Access Token](https://github.com/settings/tokens/new) with **repo** and **read:user** scopes, *setting the expiry as you see appropriate*
3. In your `go-readme-stats` repository, go to **Settings > Secrets and variables > Actions**, press **New repository secret** and add:
   - **Name:** `STAT_TOKEN`
   - **Value:** *the generated PAT* 

4. Go to the **Actions** tab in your repository, select the card generation workflow, and click **Run workflow** to generate your first SVG.
5. Copy-and-paste the link of the SVG into your README:
  ```bash
  ![Languages Card](https://raw.githubusercontent.com/<YOUR USERNAME>/go-readme-stats/output/languages-dark.svg)
  ```

### Configuration
By default, up to 6 languages are shown. If there are more, the top 5 are shown individually and the rest are grouped under *Other*.

It is recommended to configure your card by simply editing `stats.yml`, alternatively, you can use command-line flags to the run command on line 29 of `.github/workflows/generate-card.yml`.

Configuration is prioritised in this order: flags, `stats.yml`, then defaults.

```yaml
cards:
  - output: languages-dark.svg    # File name
    theme: dark
    header: Languages
    mode: geometric               # Percentage calculation method
    max_langs: 6
    ignore:
      - Jupyter Notebook
      - HTML
      - CSS
      - NSIS
      - PowerShell
      - Shell
```

The above is equivalent to:
```bash
go run ./cmd/go-readme-stats --output="languages-dark.svg" --theme=dark --header="Languages" --mode=geometric --max_langs=6 --ignore="Jupyter Notebook,HTML,CSS,NSIS,PowerShell,Shell"
```

To generate more cards, simply copy-and-paste another block ([example](https://github.com/galib-i/go-readme-stats/blob/main/stats.yml)).

### Themes
  <details>
    <summary>dark (default)</summary>
    <img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/10df753f-cb46-4448-b5bb-ebcb4f1b62a6" />
  </details>

  <details>
    <summary>soft-dark</summary>
    <img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/a2613437-4a57-4c35-ab39-a9d01214faa1" />
  </details>

  <details>
    <summary>light</summary>
    <img width="300" height="118" alt="image" src="https://github.com/user-attachments/assets/3ef24020-88b0-420b-a256-4bb2ad8c1b5e" />
  </details>

> [!TIP]
> Wrap the card with GitHub’s fragment identifiers to show cards conditionally in light or dark mode (`#gh-light-mode-only` / `#gh-dark-mode-only`):
> 
> ```markdown
> [![Languages Card](https://raw.githubusercontent.com/<YOUR USERNAME>/go-readme-stats/output/languages-dark.svg)](https://github.com/galib-i/go-readme-stats#gh-dark-mode-only)
> [![Languages Card](https://raw.githubusercontent.com/<YOUR USERNAME>/go-readme-stats/output/languages-light.svg)](https://github.com/galib-i/go-readme-stats#gh-light-mode-only)
> ```
