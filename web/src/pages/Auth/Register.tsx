import { useEffect, useRef, useState } from "react";
import {
	Mail,
	Lock,
	Eye,
	EyeOff,
	ArrowRight,
	User,
	Check,
	Github,
} from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import gsap from "gsap";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { siteConfig } from "@/content/config";
import { useRegister, useRedirectIfAuthenticated } from "@/hooks";
import { registerSchema } from "@/validations/auth";
import type { RegisterFormData } from "@/validations/auth";
import { authApi } from "@/lib/api/auth";

export default function Register() {
	const containerRef = useRef<HTMLDivElement>(null);
	const logoRef = useRef<HTMLDivElement>(null);
	const formRef = useRef<HTMLDivElement>(null);
	const [showPassword, setShowPassword] = useState(false);
	const [showConfirmPassword, setShowConfirmPassword] = useState(false);
	const [passwordStrength, setPasswordStrength] = useState(0);
	const [errorMessage, setErrorMessage] = useState("");

	const navigate = useNavigate();
	const register = useRegister();

	// Redirect if already authenticated
	useRedirectIfAuthenticated("/dashboard");

	const {
		register: registerField,
		handleSubmit,
		formState: { errors, isSubmitting },
		watch,
	} = useForm<RegisterFormData>({
		resolver: zodResolver(registerSchema),
	});

	const password = watch("password", "");
	const confirmPassword = watch("confirm_password", "");

	useEffect(() => {
		const ctx = gsap.context(() => {
			if (logoRef.current) {
				gsap.fromTo(
					logoRef.current,
					{ y: -30, opacity: 0 },
					{ y: 0, opacity: 1, duration: 0.8, ease: "power3.out" },
				);
			}

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
					},
				);
			}
		}, containerRef);

		return () => ctx.revert();
	}, []);

	useEffect(() => {
		let strength = 0;
		if (password.length >= 8) strength++;
		if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
		if (/\d/.test(password)) strength++;
		if (/[^a-zA-Z0-9]/.test(password)) strength++;
		setPasswordStrength(strength);
	}, [password]);

	const onSubmit = async (data: RegisterFormData) => {
		setErrorMessage("");
		
		try {
			await register(data);
			navigate("/dashboard");
		} catch (error: any) {
			console.error("Registration error:", error);
			setErrorMessage(
				error.response?.data?.error || "Registration failed. Please try again."
			);
		}
	};

	const handleGitHubLogin = () => {
		const githubAuthUrl = authApi.getGitHubAuthUrl();
		window.location.href = githubAuthUrl;
	};

	const getPasswordStrengthColor = () => {
		if (passwordStrength === 0) return "bg-muted";
		if (passwordStrength === 1) return "bg-red-500";
		if (passwordStrength === 2) return "bg-orange-500";
		if (passwordStrength === 3) return "bg-yellow-500";
		return "bg-green-500";
	};

	const getPasswordStrengthText = () => {
		if (passwordStrength === 0) return "";
		if (passwordStrength === 1) return "Weak";
		if (passwordStrength === 2) return "Fair";
		if (passwordStrength === 3) return "Good";
		return "Strong";
	};

	return (
		<div
			ref={containerRef}
			className="min-h-screen bg-background flex items-center justify-center relative overflow-hidden p-4"
		>
			{/* Background decorations */}
			<div className="absolute inset-0 overflow-hidden opacity-30">
				<div className="absolute top-1/4 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
				<div className="absolute bottom-1/4 left-1/4 w-96 h-96 bg-chart-2/10 rounded-full blur-3xl" />
			</div>

			<div className="w-full max-w-md relative z-10">
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
						Create your account
					</h1>
					<p className="text-muted-foreground">
						Join us and start building amazing things
					</p>
				</div>

				{/* Register Form */}
				<Card className="border border-border bg-card shadow-xl">
					<CardContent className="p-8">
						<div ref={formRef}>
							<form
								onSubmit={handleSubmit(onSubmit)}
								className="space-y-5 opacity-0"
							>
								{/* Error Message */}
								{errorMessage && (
									<div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">
										{errorMessage}
									</div>
								)}

								{/* Name Field */}
								<div>
									<label
										htmlFor="name"
										className="block text-sm font-medium text-foreground mb-2"
									>
										Full Name
									</label>
									<div className="relative">
										<User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
										<Input
											id="name"
											type="text"
											placeholder="John Doe"
											{...registerField("name")}
											className="pl-10 rounded-xl border-border focus:border-primary h-12"
										/>
									</div>
									{errors.name && (
										<p className="mt-1 text-sm text-destructive">
											{errors.name.message}
										</p>
									)}
								</div>

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
											{...registerField("email")}
											className="pl-10 rounded-xl border-border focus:border-primary h-12"
										/>
									</div>
									{errors.email && (
										<p className="mt-1 text-sm text-destructive">
											{errors.email.message}
										</p>
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
											type={
												showPassword
													? "text"
													: "password"
											}
											placeholder="Create a strong password"
											{...registerField("password")}
											className="pl-10 pr-10 rounded-xl border-border focus:border-primary h-12"
										/>
										<button
											type="button"
											onClick={() =>
												setShowPassword(!showPassword)
											}
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
										<p className="mt-1 text-sm text-destructive">
											{errors.password.message}
										</p>
									)}
									{/* Password Strength Indicator */}
									{password && (
										<div className="mt-2">
											<div className="flex gap-1 mb-1">
												{[1, 2, 3, 4].map((level) => (
													<div
														key={level}
														className={`h-1 flex-1 rounded-full transition-all duration-300 ${
															level <=
															passwordStrength
																? getPasswordStrengthColor()
																: "bg-muted"
														}`}
													/>
												))}
											</div>
											<p className="text-xs text-muted-foreground">
												Password strength:{" "}
												{getPasswordStrengthText()}
											</p>
										</div>
									)}
								</div>

								{/* Confirm Password Field */}
								<div>
									<label
										htmlFor="confirmPassword"
										className="block text-sm font-medium text-foreground mb-2"
									>
										Confirm Password
									</label>
									<div className="relative">
										<Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
										<Input
											id="confirmPassword"
											type={
												showConfirmPassword
													? "text"
													: "password"
											}
											placeholder="Confirm your password"
											{...registerField(
												"confirm_password",
											)}
											className="pl-10 pr-10 rounded-xl border-border focus:border-primary h-12"
										/>
										<button
											type="button"
											onClick={() =>
												setShowConfirmPassword(
													!showConfirmPassword,
												)
											}
											className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
										>
											{showConfirmPassword ? (
												<EyeOff className="h-5 w-5" />
											) : (
												<Eye className="h-5 w-5" />
											)}
										</button>
									</div>
									{errors.confirm_password && (
										<p className="mt-1 text-sm text-destructive">
											{errors.confirm_password.message}
										</p>
									)}
									{confirmPassword &&
										password === confirmPassword && (
											<p className="text-xs text-green-600 dark:text-green-400 mt-1 flex items-center gap-1">
												<Check className="h-3 w-3" />
												Passwords match
											</p>
										)}
								</div>

								{/* Terms and Conditions */}
								<div className="flex items-start space-x-2">
									<input
										type="checkbox"
										id="terms"
										className="w-4 h-4 mt-0.5 rounded border-border text-primary focus:ring-primary focus:ring-offset-0"
										required
									/>
									<label
										htmlFor="terms"
										className="text-sm text-muted-foreground"
									>
										I agree to the{" "}
										<Link
											to="/terms"
											className="text-primary hover:text-primary/80 transition-colors"
										>
											Terms of Service
										</Link>{" "}
										and{" "}
										<Link
											to="/privacy"
											className="text-primary hover:text-primary/80 transition-colors"
										>
											Privacy Policy
										</Link>
									</label>
								</div>

								{/* Submit Button */}
								<Button
									type="submit"
									disabled={isSubmitting}
									className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl h-12 font-semibold shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300"
								>
									{isSubmitting ? (
										"Creating account..."
									) : (
										<>
											Create Account
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
										Or sign up with
									</span>
								</div>
							</div>

							{/* Social Login */}
							<div className="space-y-3 opacity-0">
								<Button
									type="button"
									variant="outline"
									onClick={handleGitHubLogin}
									className="w-full border-border hover:bg-muted rounded-xl h-12 font-semibold transition-all duration-300 hover:scale-105"
								>
									<Github className="mr-2 h-5 w-5" />
									Continue with GitHub
								</Button>
							</div>

							{/* Sign In Link */}
							<div className="text-center mt-6 opacity-0">
								<p className="text-muted-foreground">
									Already have an account?{" "}
									<Link
										to="/login"
										className="text-primary hover:text-primary/80 font-semibold transition-colors"
									>
										Sign in
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
