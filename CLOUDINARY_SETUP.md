# Cloudinary Image Upload Integration

## Overview

The backend now supports image uploads to Cloudinary for storing AI-generated images.

## Configuration

Cloudinary credentials are stored in the `.env` file:

```env
CLOUDINARY_CLOUD_NAME=dtkxbwhzi
CLOUDINARY_API_KEY=933826486412131
CLOUDINARY_API_SECRET=0Nh4fKeSzXfxwxDoz_Zu9MBJ0FY
CLOUDINARY_UPLOAD_FOLDER=ai-of-the-world
```

## API Endpoints

### Upload Image (Protected)

**Endpoint**: `POST /api/v1/images/upload`  
**Authentication**: Required (JWT Token)  
**Content-Type**: `multipart/form-data`

**Form Fields**:
- `image` (file, required) - The image file
- `project_title` (string, required) - Title of the project
- `prompt` (string, required) - AI prompt used
- `creator_credit` (string, required) - Creator name/handle
- `technical_notes` (string, optional) - Technical details
- `model_or_tool` (string, optional) - AI model/tool used
- `tags` (string, optional) - Comma-separated tag IDs

**Example using curl**:
```bash
curl -X POST http://localhost:8080/api/v1/images/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "image=@/path/to/image.jpg" \
  -F "project_title=Cyberpunk City" \
  -F "prompt=A futuristic cyberpunk city at night" \
  -F "creator_credit=@johndoe" \
  -F "model_or_tool=Midjourney" \
  -F "tags=1,2,3"
```

**Response**:
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "project_title": "Cyberpunk City",
    "prompt": "A futuristic cyberpunk city at night",
    "image_url": "https://res.cloudinary.com/dtkxbwhzi/image/upload/v1234567890/ai-of-the-world/image.jpg",
    "status": "pending",
    "created_at": "2026-01-11T02:00:00Z"
  }
}
```

### Get All Images (Public)

**Endpoint**: `GET /api/v1/images`  
**Authentication**: Not required

**Query Parameters**:
- `status` - Filter by status (pending, approved, rejected)
- `user_id` - Filter by user ID
- `is_featured` - Filter featured images (true/false)

**Example**:
```bash
curl http://localhost:8080/api/v1/images?status=approved
```

### Get Image by ID (Public)

**Endpoint**: `GET /api/v1/images/:id`  
**Authentication**: Not required

**Example**:
```bash
curl http://localhost:8080/api/v1/images/1
```

### Delete Image (Protected)

**Endpoint**: `DELETE /api/v1/images/:id`  
**Authentication**: Required (Owner or Admin)

**Example**:
```bash
curl -X DELETE http://localhost:8080/api/v1/images/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Features

✅ **Cloudinary Integration** - Images stored in cloud  
✅ **Automatic Optimization** - Auto quality and format  
✅ **Unique Filenames** - Timestamp-based naming  
✅ **Folder Organization** - Images in `ai-of-the-world` folder  
✅ **Tag Support** - Link images to tags  
✅ **Status Workflow** - Pending → Approved/Rejected  
✅ **User Association** - Track who uploaded what  
✅ **Secure Deletion** - Remove from both DB and Cloudinary  

## Frontend Integration

### Example React/Next.js Upload Component

```javascript
import { useState } from 'react';

function ImageUpload() {
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false);

  const handleUpload = async (e) => {
    e.preventDefault();
    setUploading(true);

    const formData = new FormData();
    formData.append('image', file);
    formData.append('project_title', 'My Project');
    formData.append('prompt', 'AI prompt here');
    formData.append('creator_credit', '@username');
    formData.append('model_or_tool', 'Midjourney');
    formData.append('tags', '1,2,3');

    try {
      const token = localStorage.getItem('authToken');
      const response = await fetch('http://localhost:8080/api/v1/images/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      });

      const data = await response.json();
      if (data.success) {
        console.log('Upload successful:', data.data);
      }
    } catch (error) {
      console.error('Upload failed:', error);
    } finally {
      setUploading(false);
    }
  };

  return (
    <form onSubmit={handleUpload}>
      <input 
        type="file" 
        accept="image/*"
        onChange={(e) => setFile(e.target.files[0])}
      />
      <button type="submit" disabled={!file || uploading}>
        {uploading ? 'Uploading...' : 'Upload Image'}
      </button>
    </form>
  );
}
```

## File Size Limits

- Maximum file size: 100MB (configurable in `.env`)
- Supported formats: JPG, PNG, WEBP, GIF (static)

## Security

- ✅ Authentication required for uploads
- ✅ File type validation (images only)
- ✅ Size limit enforcement
- ✅ User ownership verification for deletion
- ✅ Admin override for deletion

## Testing

### 1. Login to get token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@aioftheworld.com", "password": "Admin@123"}'
```

### 2. Upload an image
```bash
curl -X POST http://localhost:8080/api/v1/images/upload \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "image=@test-image.jpg" \
  -F "project_title=Test Upload" \
  -F "prompt=Test prompt" \
  -F "creator_credit=@admin"
```

### 3. View uploaded images
```bash
curl http://localhost:8080/api/v1/images
```

## Troubleshooting

**Issue**: "Cloudinary not initialized"  
**Solution**: Check that your Cloudinary credentials in `.env` are correct

**Issue**: "File too large"  
**Solution**: Reduce image size or increase `MAX_UPLOAD_SIZE` in `.env`

**Issue**: "File must be an image"  
**Solution**: Ensure you're uploading a valid image file (JPG, PNG, WEBP, GIF)

---

**Cloudinary integration is ready!** 🎉 Images will be automatically uploaded to your Cloudinary account.
