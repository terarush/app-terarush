import { useState, useRef } from "react";
import { useUser } from "@/hooks";
import { profileApi } from "@/lib/api/profile";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { Camera, Loader2, Save } from "lucide-react";
import { toast } from "sonner";

export default function Profile() {
  const user = useUser();
  const [name, setName] = useState(user?.name || "");
  const [avatar, setAvatar] = useState(user?.avatar || "");
  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleAvatarClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      toast.error("Please select an image file");
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      toast.error("Image size must be less than 5MB");
      return;
    }

    setIsUploading(true);
    try {
      const response = await profileApi.uploadAvatar(file);
      setAvatar(response.avatar_url);
      toast.success("Avatar uploaded successfully");
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to upload avatar");
    } finally {
      setIsUploading(false);
    }
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await profileApi.updateProfile({
        name,
        avatar,
      });
      toast.success("Profile updated successfully");
      
      // Reload to update user context
      window.location.reload();
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to update profile");
    } finally {
      setIsSaving(false);
    }
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="flex justify-center w-full py-8 px-4">
      <div className="w-full max-w-4xl">
        <div className="mb-6">
          <h1 className="text-3xl font-bold">Profile Settings</h1>
          <p className="text-muted-foreground">Manage your account settings and profile information</p>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          {/* Avatar Section */}
          <Card className="md:col-span-1">
          <CardHeader>
            <CardTitle>Profile Picture</CardTitle>
            <CardDescription>Update your avatar</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col items-center space-y-4">
            <div className="relative group">
              <Avatar className="h-32 w-32 border-4 border-primary/20">
                <AvatarImage src={avatar ? `${import.meta.env.VITE_API_URL}${avatar}` : undefined} alt={user.name} />
                <AvatarFallback className="bg-primary/10 text-primary text-3xl font-semibold">
                  {getInitials(user.name)}
                </AvatarFallback>
              </Avatar>
              <button
                onClick={handleAvatarClick}
                disabled={isUploading}
                className="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
              >
                {isUploading ? (
                  <Loader2 className="h-8 w-8 text-white animate-spin" />
                ) : (
                  <Camera className="h-8 w-8 text-white" />
                )}
              </button>
            </div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              className="hidden"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={handleAvatarClick}
              disabled={isUploading}
            >
              <Camera className="h-4 w-4 mr-2" />
              {isUploading ? "Uploading..." : "Change Avatar"}
            </Button>
            <p className="text-xs text-muted-foreground text-center">
              JPG, PNG or GIF (max. 5MB)
            </p>
          </CardContent>
        </Card>

        {/* Profile Information */}
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>Profile Information</CardTitle>
            <CardDescription>Update your personal details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Editable Fields */}
            <div className="space-y-4">
              <div>
                <label htmlFor="name" className="text-sm font-medium">
                  Full Name
                </label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Enter your name"
                  className="mt-1.5"
                />
              </div>

              <div>
                <label htmlFor="email" className="text-sm font-medium">
                  Email Address
                </label>
                <Input
                  id="email"
                  value={user.email}
                  disabled
                  className="mt-1.5 bg-muted"
                />
                <p className="text-xs text-muted-foreground mt-1.5">
                  Email cannot be changed
                </p>
              </div>
            </div>

            <Separator />

            {/* Account Info */}
            <div className="text-sm">
              <p className="text-muted-foreground">Member Since</p>
              <p className="font-medium">
                {new Date(user.created_at).toLocaleDateString()}
              </p>
            </div>

            {/* Save Button */}
            <div className="flex justify-end pt-4">
              <Button
                onClick={handleSave}
                disabled={isSaving || !name.trim()}
              >
                {isSaving ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Saving...
                  </>
                ) : (
                  <>
                    <Save className="h-4 w-4 mr-2" />
                    Save Changes
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>
        </div>
      </div>
    </div>
  );
}
