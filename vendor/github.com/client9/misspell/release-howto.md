# Release HOWTO

since I forget.

1. Review existing tags and pick new release number

   ```bash
    git tag
   ```

2. Tag locally

   ```bash
    git tag -a v0.1.0 -m "First release"
   ```

   If things get screwed up, delete the tag with

   ```bash
   git tag -d v0.1.0
   ```

3. Test goreleaser

   TODO: how to install goreleaser

   ```bash
   ./scripts/goreleaser-dryrun.sh
   ```

4. Push

   ```bash
    git push origin v0.1.0
   ```

5. Verify release and edit notes. See [https://github.com/client9/misspell/releases](https://github.com/client9/misspell/releases)

