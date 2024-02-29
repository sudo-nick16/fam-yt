import { useEffect, useState } from "react";
import ytIcon from "./assets/youtube.png";
import searchIcon from "./assets/search.png";
import addIcon from "./assets/add.png";

const API_URL = import.meta.env.VITE_API_URL as string

type Video = {
	id: string
	title: string
	description: string
	thumbnail: string
	publishedAt: string
}

type Query = {
	id: string
	query: string
	publishedAt: string
}

function formatDate(date: string): string {
	const t = new Date(date).toLocaleString()
	return t
}

function App() {
	const [predefinedQueries, sestPredefinedQueries] = useState<Query[]>([]);
	const [query, setQuery] = useState("");
	const [videos, setVideos] = useState<Video[]>([]);
	const [loading, setLoading] = useState(false);
	const [limit, setLimit] = useState(10);
	const [page, setPage] = useState(1);

	console.log({ API_URL })

	const addSearchQuery = async (query: string) => {
		try {
			const res = await fetch(`${API_URL}/api/queries`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({ query })
			});
			const data = await res.json();
			if (data.query) {
				sestPredefinedQueries([...predefinedQueries, data.query]);
				alert("Query added successfully. Wait for a bit for the fetcher to fetch the videos.");
				return;
			}
			if (data.error) {
				alert(data.error);
			}
		} catch (err) {
			console.log("[ERROR]", err);
		}
	}

	const fetchVideos = async (query: string, limit: number, page: number) => {
		try {
			const res = await fetch(`${API_URL}/api/videos?query=${query}&limit=${limit}&pageno=${page}`, {
				method: "GET",
				headers: {
					"Content-Type": "application/json"
				}
			});
			const data = await res.json();
			if (data.videos) {
				setVideos(data.videos);
				return;
			}
			if (data.error) {
				alert(data.error);
			}
		} catch (err) {
			console.log("[ERROR]", err);
		}
	}

	const prevPage = () => {
		if (page - 1 === 0) return;
		fetchVideos(query, limit, page - 1);
		setPage(page - 1);
	}

	const nextPage = () => {
		fetchVideos(query, limit, page + 1);
		setPage(page + 1);
	}

	useEffect(() => {
		const fetchQueries = async () => {
			try {
				const res = await fetch(`${API_URL}/api/queries`, {
					method: "GET",
					headers: {
						"Content-Type": "application/json"
					}
				});
				const data = await res.json();
				if (data.queries) {
					sestPredefinedQueries(data.queries);
					setLoading(false);
					return;
				}
				if (data.error) {
					alert(data.error);
				}
				setLoading(false);
			} catch (err) {
				console.log("[ERROR]", err);
			}
		}
		fetchQueries()
	}, [])

	return (
		<div className="flex flex-col gap-4 max-w-[800px] mx-auto py-4 items-center">
			<div className="flex items-center justify-center gap-4">
				<img src={ytIcon} className="w-10 h-10" />
				<span className="text-white font-bold">Fam YT</span>
			</div>
			<div className="flex items-center gap-2 text-sm">
				<button
					onClick={() => addSearchQuery(query)}
					className="flex items-center justify-center font-semibold gap-2 p-2 bg-[#212121] rounded-md"
				>
					<img src={addIcon} className="w-4 h-4" />
					Add Query
				</button>
				<input
					type="text"
					className="outline-none font-medium bg-[#212121] text-sm px-3 py-2 rounded"
					placeholder="select or search predefined queries..."
					list="queries"
					value={query}
					onChange={(e) => setQuery(e.target.value)}
					onKeyPress={(e) => {
						if (e.key === "Enter") {
							fetchVideos(query, limit, page);
						}
					}}
				/>
				<button onClick={() => fetchVideos(query, limit, page)} className="p-2 bg-[#212121] rounded-md">
					<img src={searchIcon} className="w-5 h-5" />
				</button>
				<div className="flex gap-3 items-center">
					<label htmlFor="limit" className="font-semibold">Limit:</label>
					<input
						type="number"
						id="limit"
						className="outline-none p-2 bg-[#212121] rounded-md w-16 text-center"
						value={limit}
						onChange={(e) => setLimit(Number(e.target.value))}
					/>
				</div>
				<select
					name="sort"
					id="sort"
					className="outline-none rounded p-2 bg-[#212121]"
				>
					<option value="date">Date (Asce)</option>
					<option value="date">Date (Desc)</option>
				</select>
				<datalist id="queries">
					{
						predefinedQueries.map(query => (
							<option value={query.query} key={query.id} />
						))
					}
				</datalist>
			</div>
			<div className="">
				{
					videos.map(v => (
						<div
							key={v.id}
							className="text-sm gap-3 flex items-center p-2 bg-[#191919] border border-[#242424]"
						>
							<img src={v.thumbnail} className="w-12 h-12" />
							<div className="flex flex-col">
								<h2 className="font-medium">{v.title}</h2>
								<p
									className="line-clamp-1 font-light opacity-85 text-xs"
								>{v.description}</p>
								<span
									className="mt-1 font-light opacity-85 text-xs"
								>{formatDate(v.publishedAt)}</span>
							</div>
						</div>
					))
				}
			</div>
			<div className="flex gap-2 items-center">
				<button onClick={prevPage} className="bg-[#212121] p-1 px-4 rounded">&lt;</button>
				<span className="bg-[#212121] p-1 px-4 rounded w-12 flex items-center justify-center">{page}</span>
				<button onClick={nextPage} className="bg-[#212121] p-1 px-4 rounded">&gt;</button>
			</div>
		</div>
	)
}

export default App
