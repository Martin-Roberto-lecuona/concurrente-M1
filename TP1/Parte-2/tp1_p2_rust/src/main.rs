  use std::thread;
  use std::time::Instant;
  use rand::Rng;

  #[derive(Debug)]
  struct Matrix
  {
      matrix: Vec<Vec<i32>>,
      order: usize,
  }

  impl Matrix
  {
      pub fn new(order: usize) -> Self
      {
          let matrix = vec![vec![0; order]; order];
          Matrix
          {
            matrix, order
          }
      }

      pub fn fill_matrix(&mut self)
      {
          let n = self.order;
          for i in 0..n
          {
              for j in 0..n
              {
                  self.matrix[i][j] = rand::thread_rng().gen_range(0,10);
              }
          }
      }

      pub fn multiply_sequential(&self, scalar: i32) -> Self
      {
          let mut result = self.matrix.clone();
          for i in 0..self.order
          {
              for j in 0..self.order
              {
                  result[i][j] *= scalar;
              }
          }
          Matrix
          {
              matrix: result,
              order: self.order,
          }
      }

      pub fn multiply_concurrent(&self, scalar: i32, num_threads: usize) -> Self
      {
        let rows_per_thread = self.order / num_threads;
        let mut handles = Vec::with_capacity(num_threads);

        let mut result_matrix = Vec::with_capacity(self.order * self.order);
        for i in 0..num_threads
        {
            let start_row = i * rows_per_thread;
            let mut sub_matrix = self.matrix[start_row..(start_row + rows_per_thread)].to_vec();

            let handle = thread::spawn(move ||
            {
                for row in &mut sub_matrix
                {
                    for elem in row
                    {
                        *elem *= scalar;
                    }
                }
                sub_matrix
            });

            handles.push(handle);
        }

        for handle in handles
        {
            let sub_matrix = handle.join().unwrap();
            result_matrix.extend(sub_matrix);
        }

        Matrix
        {
            matrix: result_matrix,
            order: self.order,
        }
    }

  }

  fn main()
  {
      let n = 9000;
      let scalar = 7;
      let threads = 100;

      let mut original_matrix = Matrix::new(n);
      original_matrix.fill_matrix();

      let sequential_start = Instant::now();
      original_matrix.multiply_sequential(scalar);
      let sequential_duration = sequential_start.elapsed();

      let concurrent_start = Instant::now();
      original_matrix.multiply_concurrent(scalar, threads);
      let concurrent_duration = concurrent_start.elapsed();

      println!("Elapsed time for sequential multiplication: {:?} ms", sequential_duration.as_millis());
      println!("Elapsed time for concurrent multiplication: {:?} ms", concurrent_duration.as_millis());
  }