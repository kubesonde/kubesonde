import { FaDownload } from "react-icons/fa";
import { GrGraphQl } from "react-icons/gr";
import "./Home.scss";
import React, { useEffect } from "react";
import { useFilePicker } from "use-file-picker";
import { useNavigate } from "react-router-dom";
import { ProbeOutput } from "src/entities/probeOutput";

export const HomeComponent = () => {
  const navigate = useNavigate();
  const { openFilePicker: openFileSelector, filesContent, clear } = useFilePicker({
    accept: ".json",
  });

  // 👇 move navigation into an effect
  useEffect(() => {
    let isMounted = true;
  
    if (filesContent.length && isMounted) {
      const data = filesContent[0].content;
      const parsedData: ProbeOutput = JSON.parse(data);
      const filenameRaw = filesContent[0].name.replace(".json", "");
      const filename =
        filenameRaw.charAt(0).toUpperCase() + filenameRaw.slice(1);
  
      // clear file picker before navigating
      clear();
  
      navigate(`/graph/${filesContent[0].name}`, {
        state: { data: parsedData, title: filename },
      });
    }
  
    return () => {
      isMounted = false;
    };
  }, [filesContent, clear, navigate]);
  

  return (
    <div className="home">
      <header className="home__hero">
        <div className="home__badge">
          <img src="/logo257.png" alt="" className="home__logo" />
        </div>
        <h1 className="home__title">Kubesonde</h1>
        <p className="home__lede">
          Visualize your cluster's network connectivity
        </p>
      </header>

      <div className="home__blocks">
        <button
          type="button"
          className="block"
          onClick={() => openFileSelector()}
        >
          <span className="block__icon">
            <FaDownload aria-hidden="true" />
          </span>
          <span className="block__title">Upload a probe run</span>
          <span className="block__hint">
            Open a <code>.json</code> file exported by Kubesonde
          </span>
        </button>

        <button
          type="button"
          className="block"
          onClick={() => navigate("/example")}
        >
          <span className="block__icon">
            <GrGraphQl aria-hidden="true" />
          </span>
          <span className="block__title">Load the example</span>
          <span className="block__hint">
            Explore a sample probe to see how it works
          </span>
        </button>
      </div>
    </div>
  );
};
