use polars::prelude::*;
use std::env;
use std::io::Cursor;


// Tool to find mystery core levels of XPS spectra
//Table 1-1 gives the electron binding energies for the elements in their natural forms. A PDF version of this table is also available. The energies are given in electron volts relative to the vacuum level for the rare gases and for H2, N2, O2, F2, and Cl2; relative to the Fermi level for the metals; and relative to the top of the valence bands for semiconductors. Values have been taken from Ref. 1 except as follows:

// * Values taken from Ref. 2, with additional corrections

// d Values taken from Ref. 3.

// a One-particle approximation not valid owing to short core-hole lifetime.

// b Value derived from Ref. 1.

// https://xdb.lbl.gov/Section1/Table_1-1.pdf

// REFERENCES

// 1.      J. A. Bearden and A. F. Burr, “Reevaluation of X-Ray Atomic Energy Levels,” Rev. Mod. Phys. 39, 125 (1967).

// 2.      M. Cardona and L. Ley, Eds., Photoemission in Solids I: General Principles (Springer-Verlag, Berlin, 1978).

// 3.      J. C. Fuggle and N. Mårtensson, “Core-Level Binding Energies in Metals,” J. Electron Spectrosc. Relat. Phenom. 21, 275 (1980).

//Number,Element,K-1s,L1-2s,L2-2p1/2,L3-2p3/2,M1-3s,M2-3p1/2,M3-3p3/2,M4-3d3/2,M5-3d5/2,N1-4s,N2-4p1/2,N3-4p3/2,N4-4d3/2,N5-4d5/2,N6-4f5/2,N7-4f7/2,O1-5s,O2-5p1/2,O3-5p3/2,O4-5d3/2,O5-5d5/2,P1-6s,P2-6p1/2,P3-6p3/2
//dtypes = [Int32, Utf8, Float, Float, ...]
fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 3 {
        eprintln!("Usage: {} <energy> <tolerance>", args[0]);
        std::process::exit(1);
    }

    let energy: f64 = args[1].parse().unwrap_or_else(|_| {
        eprintln!("Invalid energy value: {}", args[1]);
        std::process::exit(1);
    });

    let tolerance: f64 = args[2].parse().unwrap_or_else(|_| {
        eprintln!("Invalid tolerance value: {}", args[2]);
        std::process::exit(1);
    });

    let csv_data = include_str!("cleaned_data.txt");

    let myschema = Schema::from(
        vec![
        Field::new("Number", DataType::Int32),
        Field::new("Element", DataType::Utf8),
        Field::new("K-1s", DataType::Float64),
        Field::new("L1-2s", DataType::Float64),
        Field::new("L2-2p1/2", DataType::Float64),
        Field::new("L3-2p3/2", DataType::Float64),
        Field::new("M1-3s", DataType::Float64),
        Field::new("M2-3p1/2", DataType::Float64),
        Field::new("M3-3p3/2", DataType::Float64),
        Field::new("M4-3d3/2", DataType::Float64),
        Field::new("M5-3d5/2", DataType::Float64),
        Field::new("N1-4s", DataType::Float64),
        Field::new("N2-4p1/2", DataType::Float64),
        Field::new("N3-4p3/2", DataType::Float64),
        Field::new("N4-4d3/2", DataType::Float64),
        Field::new("N5-4d5/2", DataType::Float64),
        Field::new("N6-4f5/2", DataType::Float64),
        Field::new("N7-4f7/2", DataType::Float64),
        Field::new("O1-5s", DataType::Float64),
        Field::new("O2-5p1/2", DataType::Float64),
        Field::new("O3-5p3/2", DataType::Float64),
        Field::new("O4-5d3/2", DataType::Float64),
        Field::new("O5-5d5/2", DataType::Float64),
        Field::new("P1-6s", DataType::Float64),
        Field::new("P2-6p1/2", DataType::Float64),
        Field::new("P3-6p3/2", DataType::Float64),
        ].into_iter());

    let df = CsvReader::new(Cursor::new(csv_data))
        .with_dtypes(Some(Arc::new(myschema)))
        .has_header(true)
        .finish()
        .unwrap();

    // convert the dtype of every column other than Number or Element to Float


    let newdf = filter_dataframe(df, energy, tolerance).unwrap();
    
    // println!("{:?}", newdf); // Print the dataframe
    print_results(&newdf, energy, tolerance);

}

fn filter_dataframe(df: DataFrame, energy: f64, tolerance: f64) -> Result<DataFrame,PolarsError> {
    let mut any_filter_expr = None;

    for column_name in df.get_column_names() {
        if column_name != "Number" && column_name != "Element" {
            let column = df.column(column_name)?;
            let mask1 = column.lt(energy + tolerance).unwrap();
            let mask2 = column.gt(energy - tolerance).unwrap();
            let new_filter_expr = mask1 | mask2;

            any_filter_expr = match any_filter_expr {
                Some(current_filter_expr) => Some(current_filter_expr | new_filter_expr),
                None => Some(new_filter_expr),
            };
        }
    }

    match any_filter_expr {
        Some(filter_expr) => df.filter(&filter_expr),
        None => Err(PolarsError::ComputeError("No columns to filter on".into())),
    }
}

fn print_results(df: &DataFrame, energy: f64, tolerance: f64) {
    let number_column = df.column("Number").unwrap();
    let element_column = df.column("Element").unwrap();

    for column_name in df.get_column_names() {
        if column_name != "Number" && column_name != "Element" {
            let column = df.column(column_name).unwrap();

            for i in 0..df.height() {
                let number = number_column.get(i).unwrap();
                let element = element_column.get(i).unwrap();
                let value = column.get(i).unwrap();
                
                if value > AnyValue::Float64(energy - tolerance) && value < AnyValue::Float64(energy + tolerance) {
                    println!("{} | {} | {} | \t{} eV", number, element, column_name, value);
                }
            }
        }
    }
}