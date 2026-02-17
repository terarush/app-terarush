import { useEffect, useRef, useState } from "react";
import { Mail, Lock, Eye, EyeOff, ArrowRight, Github } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import gsap from "gsap";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { siteConfig } from "@/content/config";
import { useAuth } from "@/contexts/AuthContext";
import { loginSchema } from "@/validations/auth";
import type { LoginFormData } from "@/validations/auth";

export default function Login() {
	const containerRef = useRef<HTMLDivElement>(null);
	const logoRef = useRef<HTMLDivElement>(null);
	const formRef = useRef<HTMLDivElement>(null);
	const [showPassword, setShowPassword] = useState(false);
	const [errorMessage, setErrorMessage] = useState("");
	
	const navigate = useNavigate();
	const { login, isAuthenticated } = useAuth();
	
	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<LoginFormData>({
		resolver: zodResolver(loginSchema),
	});

	useEffect(() => {
		// Redirect if already authenticated
		if (isAuthenticated) {
			navigate("/dashboard");
		}
	}, [isAuthenticated, navigate]);

	useEffect(() => {
		const ctx = gsap.context(() => {
			// Logo animation
			if (logoRef.current) {
				gsap.fromTo(
					logoRef.current,
					{ y: -30, opacity: 0 },
					{ y: 0, opacity: 1, duration: 0.8, ease: "power3.out" }
				);
			}

			// Form animation
			if (formRef.current) {
				gsap.fromTo(
					formRef.current.children,
					{ y: 30, opacity: 0 },
					{
						y: 0,
						opacity: 1,
						duration: 0.8,
						stagger: 0.1,
						delay: 0.3,
						ease: "power3.out",
					}
				);
			}
		}, containerRef);

		return () => ctx.revert();
	}, []);

	const onSubmit = async (data: LoginFormData) => {
		setErrorMessage("");
		
		try {
			await login(data);
			navigate("/dashboard");
		} catch (error: any) {
			console.error("Login error:", error);
			setErrorMessage(
				error.response?.data?.error || "Login failed. Please check your credentials."
			);
		}
	};

	return (
		<div
			ref={containerRef}
			className="min-h-screen bg-background flex items-center justify-center relative overflow-hidden p-4"
		>
			{/* Background decorations */}
			<div className="absolute inset-0 overflow-hidden opacity-30">
				<div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-chart-1/10 rounded-full blur-3xl" />
			</div>

			<div className="w-full max-w-md relative z-10">
				{/* Logo and Title */}
				<div ref={logoRef} className="text-center mb-8 opacity-0">
					<Link
						to="/"
						className="inline-flex items-center space-x-3 mb-6 group"
					>
						<img
							src="/assets/logo.png"
							alt="Terarush"
							className="h-12 w-12 rounded-lg shadow-lg group-hover:scale-110 transition-transform duration-300"
						/>
						<span className="text-3xl font-bold text-foreground">
							{siteConfig.name}
						</span>
					</Link>
					<h1 className="text-3xl font-bold text-foreground mb-2">
						Welcome back
					</h1>
					<p className="text-muted-foreground">
						Sign in to continue to your account
					</p>
				</div>

				{/* Login Form */}
				<Card className="border border-border bg-card shadow-xl">
					<CardContent className="p-8">
						<div ref={formRef}>
							<form onSubmit={handleSubmit(onSubmit)} className="space-y-6 opacity-0">
								{/* Error Message */}
								{errorMessage && (
									<div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">
										{errorMessage}
									</div>
								)}

								{/* Email Field */}
								<div>
									<label
										htmlFor="email"
										className="block text-sm font-medium text-foreground mb-2"
									>
										Email Address
									</label>
									<div className="relative">
										<Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
										<Input
											id="email"
											type="email"
											placeholder="you@example.com"
											{...register("email")}
											className="pl-10 rounded-xl border-border focus:border-primary h-12"
										/>
									</div>
									{errors.email && (
										<p className="mt-1 text-sm text-destructive">{errors.email.message}</p>
									)}
								</div>

								{/* Password Field */}
								<div>
									<label
										htmlFor="password"
										className="block text-sm font-medium text-foreground mb-2"
									>
										Password
									</label>
									<div className="relative">
										<Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
										<Input
											id="password"
											type={showPassword ? "text" : "password"}
											placeholder="Enter your password"
											{...register("password")}
											className="pl-10 pr-10 rounded-xl border-border focus:border-primary h-12"
										/>
										<button
											type="button"
											onClick={() => setShowPassword(!showPassword)}
											className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
										>
											{showPassword ? (
												<EyeOff className="h-5 w-5" />
											) : (
												<Eye className="h-5 w-5" />
											)}
										</button>
									</div>
									{errors.password && (
										<p className="mt-1 text-sm text-destructive">{errors.password.message}</p>
									)}
								</div>

								{/* Remember Me & Forgot Password */}
								<div className="flex items-center justify-between">
									<label className="flex items-center space-x-2 cursor-pointer">
										<input
											type="checkbox"
											className="w-4 h-4 rounded border-border text-primary focus:ring-primary focus:ring-offset-0"
										/>
										<span className="text-sm text-muted-foreground">
											Remember me
										</span>
									</label>
									<Link
										to="/forgot-password"
										className="text-sm text-primary hover:text-primary/80 transition-colors"
									>
										Forgot password?
									</Link>
								</div>

								{/* Submit Button */}
								<Button
									type="submit"
									disabled={isSubmitting}
									className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl h-12 font-semibold shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300"
								>
									{isSubmitting ? (
										"Signing in..."
									) : (
										<>
											Sign In
											<ArrowRight className="ml-2 h-5 w-5" />
										</>
									)}
								</Button>
							</form>

							{/* Divider */}
							<div className="relative my-6 opacity-0">
								<div className="absolute inset-0 flex items-center">
									<div className="w-full border-t border-border"></div>
								</div>
								<div className="relative flex justify-center text-sm">
									<span className="px-4 bg-card text-muted-foreground">
										Or continue with
									</span>
								</div>
							</div>

							{/* Social Login */}
							<div className="space-y-3 opacity-0">
								<Button
									type="button"
									variant="outline"
									className="w-full border-border hover:bg-muted rounded-xl h-12 font-semibold transition-all duration-300 hover:scale-105"
								>
									<Github className="mr-2 h-5 w-5" />
									Continue with GitHub
								</Button>
							</div>

							{/* Sign Up Link */}
							<div className="text-center mt-6 opacity-0">
								<p className="text-muted-foreground">
									Don't have an account?{" "}
									<Link
										to="/register"
										className="text-primary hover:text-primary/80 font-semibold transition-colors"
									>
										Sign up
									</Link>
								</p>
							</div>
						</div>
					</CardContent>
				</Card>

				{/* Back to Home */}
				<div className="text-center mt-6">
					<Link
						to="/"
						className="inline-flex items-center text-muted-foreground hover:text-foreground transition-colors"
					>
						← Back to home
					</Link>
				</div>
			</div>
		</div>
	);
}
