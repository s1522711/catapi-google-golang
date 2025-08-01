# Pterodactyl Setup Guide

## 🚀 **Deploying Google Photos Album API on Pterodactyl**

### **Step 1: Import the Egg**

1. **Access your Pterodactyl Admin Panel**
2. **Go to Nests** → **Create New Nest** (if needed)
3. **Go to Eggs** → **Import Egg**
4. **Upload the `egg-google-photos-api.json` file**

### **Step 2: Create a New Server**

1. **Go to Servers** → **Create New Server**
2. **Select the "Google Photos Album API" egg**
3. **Configure server settings:**
   - **Name**: `Google Photos API`
   - **Description**: `Random image API from Google Photos album`
   - **Memory**: `512 MB` (minimum)
   - **CPU**: `50%` (minimum)
   - **Disk**: `1 GB` (minimum)

### **Step 3: Configure Variables**

Set these variables in the server creation:

| Variable | Value | Description |
|----------|-------|-------------|
| `GITHUB_REPO` | `https://github.com/s1522711/catapi-google-golang.git` | GitHub repository URL (default) |
| `ALBUM_URL` | `https://photos.app.goo.gl/YOUR_ALBUM_ID` | Your public Google Photos album URL |
| `PORT` | `8080` | Port for the API (usually 8080) |

### **Step 4: Automatic GitHub Integration**

The egg will automatically:
- ✅ **Clone the latest version** from GitHub on first install
- ✅ **Pull latest changes** every time the server starts
- ✅ **Install dependencies** automatically
- ✅ **Build the application** ready to run

**No manual file upload required!** The egg handles everything automatically and always runs the latest version.

### **Step 5: Start the Server**

1. **Click "Start Server"**
2. **Monitor the console for:**
   ```
   # First time:
   Cloning into '.'...
   remote: Enumerating objects...
   
   # Every startup:
   Already up to date.
   # or
   Updating X..Y
   go: downloading github.com/gin-gonic/gin
   Using environment variables - Album URL: https://photos.app.goo.gl/...
   Auto-refresh started: cache will refresh every 12 hours
   Server starting on port 8080
   ```

### **Step 6: Access Your API**

**Web Interface:**
```
http://your-server-ip:8080
```

**API Endpoints:**
```
http://your-server-ip:8080/api/random
http://your-server-ip:8080/api/img.png
http://your-server-ip:8080/api/images
http://your-server-ip:8080/api/refresh
```

## 🔧 **Configuration Options**

### **Environment Variables (Recommended)**
- `ALBUM_URL` - Google Photos album URL
- `PORT` - Server port (default: 8080)

### **Config File (Alternative)**
Create `config.json`:
```json
{
  "album_url": "https://photos.app.goo.gl/YOUR_ALBUM_ID",
  "port": 8080
}
```

## 📊 **Resource Requirements**

**Minimum:**
- **Memory**: 512 MB
- **CPU**: 50%
- **Disk**: 1 GB

**Recommended:**
- **Memory**: 1 GB
- **CPU**: 100%
- **Disk**: 2 GB

## 🔍 **Troubleshooting**

### **Common Issues:**

1. **"Album URL not found"**
   - Check `ALBUM_URL` environment variable
   - Verify the Google Photos album is public

2. **"Port already in use"**
   - Change the `PORT` variable
   - Check if another service is using port 8080

3. **"No images found in album"**
   - Verify the album URL is correct
   - Ensure the album contains images
   - Check if the album is publicly accessible

4. **Build errors**
   - Ensure all files are uploaded
   - Check `go.mod` and `go.sum` are present

### **Logs to Monitor:**
```
Refreshing image cache...
Cached X images
Auto-refresh started: cache will refresh every 12 hours
Server starting on port 8080
```

## 🌐 **External Access**

To make your API accessible from the internet:

1. **Configure your firewall** to allow port 8080
2. **Set up a reverse proxy** (nginx/Apache) if needed
3. **Use a domain name** for easier access

## 🔄 **Automatic Updates**

The API automatically updates every time you start the server:

1. **Stop the server** in Pterodactyl
2. **Start the server** - automatically pulls latest changes
3. **No reinstall needed!** Just restart to get updates

The egg will automatically:
- ✅ **Pull latest changes** from GitHub on every startup
- ✅ **Install updated dependencies** if needed
- ✅ **Always run the latest version**

## 📈 **Scaling**

For high-traffic deployments:

1. **Increase memory** to 2-4 GB
2. **Add CPU cores** for better performance
3. **Consider load balancing** for multiple instances
4. **Monitor cache refresh** performance

---

**🎉 Your Google Photos Album API is now ready to serve random images!** 