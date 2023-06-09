import type { NextPage } from 'next';
import Head from 'next/head';
import axios from 'axios';
import { FormEventHandler, useEffect, useState } from 'react';

interface Task {
	id: number;
	task_name: string;
	owner_name: string;
	status: boolean;
}

const Home: NextPage = () => {
	const [allTasks, setAllTasks] = useState<Task[]>([]);
	const [task, setTask] = useState('');
	const [owner, setOwner] = useState('');
	const [status, setStatus] = useState(false);

	const handleSubmit: FormEventHandler<HTMLFormElement> = async e => {
		e.preventDefault();

		try {
			let { data } = await axios.post('http://localhost:4000', {
				task_name: task,
				owner_name: owner,
			});

			console.log(data);
			window.location.reload();
		} catch (error) {
			console.log(error);
		}
	};

	useEffect(() => {
		const getData = async () => {
			let { data } = await axios.get('http://localhost:4000');

			setAllTasks(data);
		};

		getData();
	}, []);

	return (
		<div>
			<Head>
				<title>Go Fiber + MySQL CRUD</title>
				<meta name='description' content='Generated by create next app' />
				<link rel='icon' href='/favicon.ico' />
			</Head>

			<main className='flex flex-col items-center pt-8'>
				<h1 className='mb-8'>Go Todo App</h1>
				<form onSubmit={handleSubmit}>
					<div className='my-2'>
						<label className='block' htmlFor='task_name'>
							Task:
						</label>
						<input
							className='px-2 py-1'
							type='text'
							name='task'
							value={task}
							onChange={e => setTask(e.target.value)}
              required
						/>
					</div>
					<div className='my-2'>
						<label className='block' htmlFor='owner_name'>
							Owner:
						</label>
						<input
							className='px-2 py-1'
							type='text'
							name='task'
							value={owner}
							onChange={e => setOwner(e.target.value)}
              required
						/>
					</div>
					{/* <div className='pt-2'>
						<label className='mr-4' htmlFor='status'>
							Status:
						</label>
						<select
							name='status'
							onChange={e => {
								if (e.target.value == 'true') {
									setStatus(true);
								} else {
									setStatus(false);
								}
							}}
						>
							<option value='false'>incomplete</option>
							<option value='true'>completed</option>
						</select>
					</div> */}
					<div className='flex justify-center mt-4'>
						<button
							className='border border-white p-2 hover:cursor-pointer hover:bg-slate-100 hover:text-black'
							type='submit'
						>
							Add Task
						</button>
					</div>
				</form>

				<section className='border border-white mt-8 p-4'>
					{allTasks.map((task: Task) => (
						<div key={task.id} className='border border-white p-2'>
							<h4>Task:</h4>
							<p className='mb-4'>{task.task_name}</p>
							<h4>Owner:</h4>
							<p className='mb-4'>{task.owner_name}</p>
							<h4>Completed:</h4>
							<p>{task.status ? 'yes' : 'no'}</p>
						</div>
					))}
				</section>
			</main>
		</div>
	);
};

export default Home;
