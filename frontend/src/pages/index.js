import { useEffect, useRef, useState } from "react";

export default function GridApp() {
  const [inputGridSize, setInputGridSize] = useState("");
  const [gridSize, setGridSize] = useState(null);
  const [provideInitialState, setProvideInitialState] = useState(false);
  const [initialSubmitted, setInitialSubmitted] = useState(false);

  const [userPoints, setUserPoints] = useState([]);

  const [pointSets, setPointSets] = useState([]);
  const [currentSetIndex, setCurrentSetIndex] = useState(0);
  const [autoCycle, setAutoCycle] = useState(true);

  const canvasRef = useRef(null);
  const inputCanvasRef = useRef(null);

  useEffect(() => {
    if (!gridSize || !initialSubmitted) return;

    const bodyData = {
      ...(provideInitialState && { points: userPoints }),
    };

    fetch(`http://localhost:8080/api/points?gridSize=${gridSize}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(bodyData),
    })
      .then((res) => res.json())
      .then((data) => {
        setPointSets(data);
        setCurrentSetIndex(0);
      })
      .catch((err) => console.error("Error fetching points:", err));
  }, [gridSize, initialSubmitted, provideInitialState, userPoints]);

  useEffect(() => {
    if (!gridSize || pointSets.length === 0 || !canvasRef.current) return;

    const cellSize = 40;
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    canvas.width = gridSize * cellSize;
    canvas.height = gridSize * cellSize;

    function drawGrid(pointsData) {
      const points = pointsData.Points;
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      ctx.strokeStyle = "#ddd";
      for (let i = 0; i <= gridSize; i++) {
        ctx.beginPath();
        ctx.moveTo(i * cellSize, 0);
        ctx.lineTo(i * cellSize, canvas.height);
        ctx.moveTo(0, i * cellSize);
        ctx.lineTo(canvas.width, i * cellSize);
        ctx.stroke();
      }

      ctx.fillStyle = "blue";
      if (points != null) {
        points.forEach((point) => {
          const adjustedY = gridSize - 1 - point.Y;
          ctx.fillRect(point.X * cellSize, adjustedY * cellSize, cellSize, cellSize);
        });
      }
    }

    if (pointSets[currentSetIndex]) {
      drawGrid(pointSets[currentSetIndex]);
    }
  }, [gridSize, pointSets, currentSetIndex]);

  useEffect(() => {
    if (pointSets.length === 0 || !autoCycle) return;

    const interval = setInterval(() => {
      setCurrentSetIndex((prev) => (prev + 1) % pointSets.length);
    }, 200);

    return () => clearInterval(interval);
  }, [pointSets.length, autoCycle]);

  const handlePrevious = () => {
    setAutoCycle(false);
    setCurrentSetIndex((prev) =>
      prev === 0 ? pointSets.length - 1 : prev - 1
    );
  };

  const handleNext = () => {
    setAutoCycle(false);
    setCurrentSetIndex((prev) => (prev + 1) % pointSets.length);
  };

  const toggleAutoCycle = () => {
    setAutoCycle((prev) => !prev);
  };

  const handleGridSizeSubmit = (e) => {
    e.preventDefault();
    const size = Number(inputGridSize);
    if (size > 0) {
      setGridSize(size);
      if (!provideInitialState) {
        setInitialSubmitted(true);
      }
    }
  };

  useEffect(() => {
    if (!gridSize || !provideInitialState || !inputCanvasRef.current) return;

    const cellSize = 40;
    const canvas = inputCanvasRef.current;
    const ctx = canvas.getContext("2d");
    canvas.width = gridSize * cellSize;
    canvas.height = gridSize * cellSize;

    function drawInitialGrid() {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.strokeStyle = "#ddd";
      for (let i = 0; i <= gridSize; i++) {
        ctx.beginPath();
        ctx.moveTo(i * cellSize, 0);
        ctx.lineTo(i * cellSize, canvas.height);
        ctx.moveTo(0, i * cellSize);
        ctx.lineTo(canvas.width, i * cellSize);
        ctx.stroke();
      }
      ctx.fillStyle = "red";
      userPoints.forEach((pt) => {
        const adjustedY = gridSize - 1 - pt.y;
        ctx.fillRect(pt.x * cellSize, adjustedY * cellSize, cellSize, cellSize);
      });
    }

    drawInitialGrid();

    const handleClick = (e) => {
      const rect = canvas.getBoundingClientRect();
      const clickX = e.clientX - rect.left;
      const clickY = e.clientY - rect.top;
      const cellX = Math.floor(clickX / cellSize);
      const cellY = gridSize - 1 - Math.floor(clickY / cellSize);

      setUserPoints((prev) => {
        if (prev.some((pt) => pt.x === cellX && pt.y === cellY)) return prev;
        return [...prev, { x: cellX, y: cellY }];
      });
    };

    canvas.addEventListener("click", handleClick);
    return () => canvas.removeEventListener("click", handleClick);
  }, [gridSize, provideInitialState, userPoints]);

  const handleInitialStateSubmit = () => {
    setInitialSubmitted(true);
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
      <h1 className="text-2xl font-bold mb-4">John Conway's Game of Life</h1>
      {!gridSize && (
        <form onSubmit={handleGridSizeSubmit} className="mb-4 flex flex-col items-center">
          <div className="mb-2">
            <label htmlFor="gridSize" className="mr-2">
              Enter Grid Size:
            </label>
            <input
              id="gridSize"
              type="number"
              value={inputGridSize}
              onChange={(e) => setInputGridSize(e.target.value)}
              className="border border-gray-400 p-1 rounded w-20"
              min="1"
              required
            />
          </div>
          <div className="mb-2">
            <span className="mr-2">Provide initial grid state?</span>
            <label className="mr-2">
              <input
                type="radio"
                name="initialState"
                value="yes"
                onChange={() => setProvideInitialState(true)}
                required
              />{" "}
              Yes
            </label>
            <label>
              <input
                type="radio"
                name="initialState"
                value="no"
                onChange={() => setProvideInitialState(false)}
              />{" "}
              No
            </label>
          </div>
          <button
            type="submit"
            className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
          >
            Set Grid Size
          </button>
        </form>
      )}

      {gridSize && provideInitialState && !initialSubmitted && (
        <div className="mb-4 flex flex-col items-center">
          <p className="mb-2">Click on the grid to add your points.</p>
          <canvas
            ref={inputCanvasRef}
            className="border border-gray-400 mb-2"
          />
          <p className="mb-2">
            {userPoints.length > 0
              ? `You have selected ${userPoints.length} point(s).`
              : "No points selected yet."}
          </p>
          <button
            onClick={handleInitialStateSubmit}
            className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
          >
            Submit Initial Grid State
          </button>
        </div>
      )}

      {gridSize && initialSubmitted && (
        <>
          <canvas ref={canvasRef} className="border border-gray-400 mb-4" />
          <div className="flex space-x-4">
            <button
              onClick={handlePrevious}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Previous
            </button>
            <button
              onClick={handleNext}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Next
            </button>
            <button
              onClick={toggleAutoCycle}
              className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
            >
              {autoCycle ? "Stop Auto Cycle" : "Start Auto Cycle"}
            </button>
          </div>
        </>
      )}
    </div>
  );
}
