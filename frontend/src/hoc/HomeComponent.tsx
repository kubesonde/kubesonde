import { FaDownload } from "react-icons/fa";
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

      <button
        type="button"
        className="dropzone"
        onClick={() => openFileSelector()}
      >
        <span className="dropzone__icon">
          <FaDownload aria-hidden="true" />
        </span>
        <span className="dropzone__title">Open a probe run</span>
        <span className="dropzone__hint">
          Choose a probe output file to visualize its network graph
        </span>
        <span className="dropzone__cta">Choose file</span>
        <span className="dropzone__note">
          Accepts <code>.json</code> exported by Kubesonde
        </span>
      </button>
    </div>
  );
};
