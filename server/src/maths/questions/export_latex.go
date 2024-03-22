package questions

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

const latexHeader = `
% ------------------- HEADER START ---------------------------- %
% Required packages 
\usepackage{amsmath}
\usepackage[inline]{enumitem}
\usepackage{amssymb}
\usepackage[table]{xcolor}
\usepackage{tikz}
\usepackage{environ}

\definecolor{isyroPropColor}{gray}{0.9}

% Custom commands and settings used by Isyro
\newcommand{\isyroFieldHeight}{\phantom{$\sum_{\sum}^{\sum}$}} %% field height

\newcommand{\isyroNumberField}{
	\fbox{\begin{minipage}{1cm}\isyroFieldHeight \end{minipage}}
}

\newcommand{\isyroExpressionField}[1]{
	\fbox{\begin{minipage}{#1 cm}\isyroFieldHeight \\ \isyroFieldHeight \end{minipage}}
}

\newcommand{\isyroQCMSquare}{\raisebox{-.25\height}{\huge$\square$}}

\newcommand{\R}{\mathbb{R}}
\newcommand{\Q}{\mathbb{Q}}
\newcommand{\D}{\mathbb{D}}
\newcommand{\Z}{\mathbb{Z}}
\newcommand{\N}{\mathbb{N}}

% Hack to adjust the figure to the page width
% See https://tex.stackexchange.com/questions/6388/how-to-scale-a-tikzpicture-to-textwidth
\makeatletter
\newsavebox{\measure@tikzpicture}
\NewEnviron{scaletikzpicturetowidth}[1]{%
	\def\tikz@width{#1}%
	\def\tikzscale{1}\begin{lrbox}{\measure@tikzpicture}%
		\BODY
	\end{lrbox}%
	\pgfmathparse{#1/\wd\measure@tikzpicture}%
	\edef\tikzscale{\pgfmathresult}%
	\BODY
}
\makeatother

% ------------------- HEADER END ---------------------------- %
`

func (qu EnonceInstance) toLatex(standalone bool) string {
	chunks := make([]string, len(qu))
	// we add an extra new line between two text blocks
	isPreviousText := false
	for i, p := range qu {
		_, isText := p.(TextInstance)

		if isPreviousText && isText {
			chunks[i] = "\n \\vspace{0.3cm} \n" + p.toLatex()
		} else {
			chunks[i] = p.toLatex()
		}

		isPreviousText = isText
	}
	// we add line breaks for clarity, they are ignored by the latex compiler anyway
	// we also disable indent inside one question
	code := fmt.Sprintf(`
	%% -------------- Question --------------
	{ \setlength{\parindent}{0pt}
	%s
	}`, strings.Join(chunks, "\n"))

	if standalone { // add custom defs
		return latexHeader + code
	}
	return code
}

func (qu EnonceInstance) ToLatex() string {
	return qu.toLatex(true)
}

func InstancesToLatex(questions []EnonceInstance) string {
	latexCodes := make([]string, len(questions))
	for i, qu := range questions {
		latexCodes[i] = `\item ` + qu.toLatex(false)
	}
	return fmt.Sprintf(`%s 
	
	\begin{enumerate}
	%s
	\end{enumerate}`, latexHeader, strings.Join(latexCodes, "\n"))
}

func textOrMathToLatex(c client.TextOrMath) string {
	if c.IsMath {
		return "$ " + c.Text + " $"
	}
	return c.Text
}

func lineToLatexCode(line client.TextLine) string {
	chunks := make([]string, len(line))
	for i, c := range line {
		chunks[i] = textOrMathToLatex(c)
	}
	out := strings.ReplaceAll(strings.Join(chunks, ""), "\n", `~\\`)
	return out
}

func (ti TextInstance) toLatex() string {
	text := lineToLatexCode(ti.Parts)

	attrs := ""
	if ti.Bold {
		attrs += `\bfseries `
	}
	if ti.Italic {
		attrs += `\itshape `
	}
	if ti.Smaller {
		attrs += `\small `
	}

	return fmt.Sprintf("{%s%s}", attrs, text)
}

func (fi FormulaDisplayInstance) toLatex() string {
	return "$$" + strings.Join(fi, " ") + "$$"
}

func tableRow(row []client.TextOrMath) string {
	cells := make([]string, len(row))
	for j, cell := range row {
		cells[j] = textOrMathToLatex(cell)
	}
	return strings.Join(cells, " & ") + `\\` + "\n"
}

func (ti TableInstance) toLatex() string {
	if len(ti.Values) == 0 || len(ti.Values[0]) == 0 {
		return ""
	}
	hasVertHeaders := len(ti.VerticalHeaders) != 0
	nbRows, nbCols := len(ti.Values), len(ti.Values[0]) // without any header

	if hasVertHeaders {
		nbCols += 1
	}

	header := ""
	if len(ti.HorizontalHeaders) != 0 {
		row := ti.HorizontalHeaders
		if hasVertHeaders {
			row = append([]client.TextOrMath{{}}, row...)
		}

		header = `\hline
		\rowcolor{pink} ` + tableRow(row)
	}

	colDeclaration := "|" + strings.Repeat(" c |", nbCols)

	rows := make([]string, nbRows)
	for i, row := range ti.Values {
		rowCode := tableRow(row)
		if hasVertHeaders {
			rowCode = fmt.Sprintf(`\cellcolor{cyan} %s & %s`, textOrMathToLatex(ti.VerticalHeaders[i]), rowCode)
		}
		rows[i] = rowCode
	}

	return fmt.Sprintf(`
	\begin{center}
		\begin{tabular}{%s}
			%s
			\hline
			%s
			\hline
		\end{tabular}
	\end{center}
`, colDeclaration, header, strings.Join(rows, `\hline`+"\n"))
}

func (vi VariationTableInstance) toLatex() string { return "TODO" }
func (si SignTableInstance) toLatex() string      { return "TODO" }

// return color and opacity
func tikzColorArg(c repere.ColorHex) (string, string) {
	a, r, g, b := c.ToARGB()

	color := fmt.Sprintf("{rgb,255:red,%d;green,%d;blue,%d}", r, g, b)
	opacity := fmt.Sprintf("%.02f", float64(a)/255.)
	return color, opacity
}

func (fi FigureInstance) toLatex() string {
	dr := fi.Figure.Drawings
	ox, oy := fi.Figure.Bounds.Origin.X, fi.Figure.Bounds.Origin.Y // all points must be translated

	gridColor := "gray"
	if !fi.Figure.ShowGrid {
		gridColor = "gray!0"
	}
	origin := ""
	if fi.Figure.ShowOrigin {
		origin = fmt.Sprintf(`\filldraw[black] (%.02f,%.02f) circle (3pt) node[anchor=south west] {$O$}; %% origin`,
			fi.Figure.Bounds.Origin.X, fi.Figure.Bounds.Origin.Y)
	}

	var drawings []string

	for name, point := range dr.Points {
		color, opacity := tikzColorArg(point.Color)
		code := fmt.Sprintf(`
			\coordinate (%s) at (%.02f,%.02f);
			\filldraw[color=%s, opacity=%s] (%s) circle (3pt) node[anchor=south west] {$%s$};`,
			name, point.Point.Point.X+ox, point.Point.Point.Y+oy,
			color, opacity, name, name)
		drawings = append(drawings, code)
	}

	lines := dr.Lines

	for _, segment := range dr.Segments {
		color, opacity := tikzColorArg(segment.Color)
		from, to := dr.Points[segment.From].Point.Point, dr.Points[segment.To].Point.Point
		kind := ""
		if segment.Kind == repere.SKVector {
			kind = ", ->"
		} else if segment.Kind == repere.SKLine {
			// infer the affine line and draw it later
			a, b := repere.InferLine(from, to)
			line := repere.Line{A: a, B: b, Label: segment.LabelName, Color: segment.Color}
			lines = append(lines, line)
			continue
		}
		code := fmt.Sprintf(`\draw[color=%s, opacity=%s %s] (%.02f,%.02f) -- (%.02f,%.02f) node[anchor=south west] {$%s$};`,
			color, opacity, kind, from.X+ox, from.Y+oy, to.X+ox, to.Y+oy, segment.LabelName)
		drawings = append(drawings, code)
	}

	for _, line := range lines {
		color, opacity := tikzColorArg(line.Color)
		from, to := line.Bounds(fi.Figure.Bounds)
		code := fmt.Sprintf(`\draw[color=%s, opacity=%s] (%.02f,%.02f) -- (%.02f,%.02f) node[anchor=south west] {$%s$};`,
			color, opacity, from.X+ox, from.Y+oy, to.X+ox, to.Y+oy, line.Label)
		drawings = append(drawings, code)
	}

	for _, circle := range dr.Circles {
		color, opacity := tikzColorArg(circle.LineColor)
		fillColor, fillOpacity := tikzColorArg(circle.FillColor)

		code := fmt.Sprintf(`\filldraw[color=%s, draw opacity=%s, fill=%s, fill opacity=%s] (%.02f,%.02f) circle (%.02f) node[anchor=south west] {$%s$};`,
			color, opacity, fillColor, fillOpacity,
			circle.Center.X+ox, circle.Center.Y+oy, circle.Radius, circle.Legend)
		drawings = append(drawings, code)
	}

	for _, area := range dr.Areas {
		points := make([]string, len(area.Points))
		for i, p := range area.Points {
			points[i] = fmt.Sprintf("(%s)", p)
		}
		path := strings.Join(points, " -- ")
		color, opacity := tikzColorArg(area.Color)
		code := fmt.Sprintf(`\path[fill=%s, fill opacity=%s] %s;`, color, opacity, path)
		drawings = append(drawings, code)
	}

	return fmt.Sprintf(`
	\begin{scaletikzpicturetowidth}{\textwidth}
	\begin{tikzpicture}[scale=\tikzscale]
		\draw[%[1]s] (-0.25,-0.25) grid (%[2]d.25,%[3]d.25);
		\draw[thick, black, ->] (-0.25, %.02[5]f) -- (%[2]d.25,%[5]f) node[anchor=south west] {$x$};
		\draw[thick, black, ->] (%.02[4]f, -0.25) -- (%.02[4]f, %[3]d.25) node[anchor=south west] {$y$};
		%[6]s

		%s
	\end{tikzpicture}
	\end{scaletikzpicturetowidth}
	`, gridColor,
		fi.Figure.Bounds.Width, fi.Figure.Bounds.Height,
		fi.Figure.Bounds.Origin.X, fi.Figure.Bounds.Origin.Y,
		origin,
		strings.Join(drawings, "\n"),
	)
}

func (fi FunctionsGraphInstance) toLatex() string { return "TODO" }

func (ni NumberFieldInstance) toLatex() string { return `\isyroNumberField` }

func (ei ExpressionFieldInstance) toLatex() string {
	width := float64(ei.sizeHint())
	// map from 1 - 30 to 2cm - 10cm
	cm := 2. + (10-2)*(width-1)/(30-1)
	if ei.LabelLaTeX != "" { // add a new line and display mode
		return fmt.Sprintf(`$$ %s \isyroExpressionField{%.2f}$$`, ei.LabelLaTeX, cm)
	}
	return fmt.Sprintf(`\isyroExpressionField{%.2f}`, cm)
}

// requires the following latex packages
//   - \usepackage[inline]{enumitem}
//   - \usepackage{amssymb}
func (ri RadioFieldInstance) toLatex() string {
	props := ri.proposals()
	choices := make([]string, len(props))
	for i, p := range props {
		choices[i] = `\item ` + lineToLatexCode(p)
	}

	return fmt.Sprintf(`\begin{itemize}[label={\isyroQCMSquare}] %% vertical align
    %s
	\end{itemize}
	`, strings.Join(choices, "\n"))
}

func (di DropDownFieldInstance) toLatex() string {
	props := RadioFieldInstance(di).proposals()
	choices := make([]string, len(props))
	for i, p := range props {
		choices[i] = `\item ` + lineToLatexCode(p)
	}
	return fmt.Sprintf(`\begin{itemize*}[label={\isyroQCMSquare}] %% vertical align
    %s
	\end{itemize*}
	`, strings.Join(choices, "\n"))
}

func (oi OrderedListFieldInstance) toLatex() string {
	props := oi.proposals()
	choices := make([]string, len(props))
	for i, p := range props {
		choices[i] = fmt.Sprintf(`\colorbox{isyroPropColor}{%s}`, lineToLatexCode(p))
	}
	box := fmt.Sprintf(`\isyroExpressionField{13}`)
	label := ""
	if oi.Label != "" {
		label = "$" + oi.Label + "$"
	}
	return fmt.Sprintf(`~\\ \begin{center} %s %s \\ 
	\vspace{0.5em}
	\textit{\small \'Eléments à ordonner:} %s
	\end{center}
	`, label, box, strings.Join(choices, " , "))
}

func (pi VectorFieldInstance) toLatex() string {
	if pi.DisplayColumn {
		return `$\begin{pmatrix} \isyroNumberField \\ \isyroNumberField \end{pmatrix}$`
	} else {
		return `(\isyroNumberField ; \isyroNumberField)`
	}
}

func (fi GeometricConstructionFieldInstance) toLatex() string { return "TODO" }
func (vi VariationTableFieldInstance) toLatex() string        { return "TODO" }
func (si SignTableFieldInstance) toLatex() string             { return "TODO" }
func (fi FunctionPointsFieldInstance) toLatex() string        { return "TODO" }
func (ti TreeInstance) toLatex() string                       { return "TODO" }
func (ti TreeFieldInstance) toLatex() string                  { return "TODO" }
func (pi ProofFieldInstance) toLatex() string                 { return "TODO" }
func (pi TableFieldInstance) toLatex() string                 { return "TODO" }
