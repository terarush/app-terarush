import Navbar from "@/components/navbar";
import Hero from "@/components/hero";
import Services from "@/components/services";
import Products from "@/components/products";
import OpenSource from "@/components/opensource";
import About from "@/components/about";
import WhyUs from "@/components/whyus";
import Contact from "@/components/contact";
import Footer from "@/components/footer";
import { useSmoothScroll } from "@/hooks/useSmoothScroll";

export default function Index() {
	useSmoothScroll();
	return (
		<>
			<Navbar />
			<Hero />
			<Services />
			<Products />
			<OpenSource />
			<About />
			<WhyUs />
			<Contact />
			<Footer />
		</>
	);
}
